package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/nczz/phishguard/config"
	"github.com/nczz/phishguard/internal/crypto"
	"github.com/nczz/phishguard/internal/db"
	"github.com/nczz/phishguard/internal/mailer"
	"github.com/nczz/phishguard/internal/model"
	"github.com/nczz/phishguard/internal/repo"
	"github.com/nczz/phishguard/internal/service"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// rateLimiters holds per-provider singleton rate limiters (persists across polling cycles)
var rateLimiters = make(map[string]*mailer.RateLimiter)

func getRateLimiter(providerType string) *mailer.RateLimiter {
	if rl, ok := rateLimiters[providerType]; ok {
		return rl
	}
	rl := mailer.NewRateLimiter(providerType)
	rateLimiters[providerType] = rl
	return rl
}

func main() {
	cfg := config.Load()
	database := db.Init(cfg.DBDSN)

	resultRepo := &repo.ResultRepo{DB: database}
	campaignRepo := &repo.CampaignRepo{DB: database}
	recipientRepo := &repo.RecipientRepo{DB: database}
	scenarioRepo := &repo.ScenarioRepo{DB: database}

	campaignSvc := &service.CampaignService{
		CampaignRepo: campaignRepo, ResultRepo: resultRepo,
		RecipientRepo: recipientRepo, ScenarioRepo: scenarioRepo,
	}

	// Start mail polling loop
	log.Println("Worker started, polling for pending emails...")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	processMailQueue(database, cfg, resultRepo, campaignRepo)
	runScheduler(database, campaignSvc)
	for range ticker.C {
		processMailQueue(database, cfg, resultRepo, campaignRepo)
		runScheduler(database, campaignSvc)
	}
}

func processMailQueue(database *gorm.DB, cfg *config.Config, resultRepo *repo.ResultRepo, campaignRepo *repo.CampaignRepo) {
	// Find campaigns in "sending" status
	campaigns, err := campaignRepo.FindByStatus(model.CampaignStatusSending)
	if err != nil {
		log.Printf("error finding sending campaigns: %v", err)
		return
	}

	for _, c := range campaigns {
		processCampaign(database, cfg, resultRepo, campaignRepo, &c)
	}
}

func processCampaign(database *gorm.DB, cfg *config.Config, resultRepo *repo.ResultRepo, campaignRepo *repo.CampaignRepo, campaign *model.Campaign) {
	// Get SMTP profile
	var smtp model.SMTPProfile
	if err := database.First(&smtp, campaign.SMTPProfileID).Error; err != nil {
		log.Printf("campaign %d: smtp profile not found: %v", campaign.ID, err)
		return
	}

	// Get template
	var tmpl model.EmailTemplate
	templateID := campaign.TemplateID
	if templateID == nil && campaign.ScenarioID != nil {
		var scenario model.Scenario
		if err := database.First(&scenario, *campaign.ScenarioID).Error; err == nil {
			templateID = scenario.TemplateID
		}
	}
	if templateID == nil {
		log.Printf("campaign %d: no template found", campaign.ID)
		return
	}
	if err := database.First(&tmpl, *templateID).Error; err != nil {
		log.Printf("campaign %d: template not found: %v", campaign.ID, err)
		return
	}

	// Create mailer
	m, err := buildMailer(&smtp, cfg.EncryptKey)
	if err != nil {
		log.Printf("campaign %d: failed to create mailer: %v", campaign.ID, err)
		return
	}

	// Claim results: SELECT FOR UPDATE SKIP LOCKED → mark as "sending"
	var results []model.Result
	now := time.Now().UTC()
	err = database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("campaign_id = ? AND status = ? AND (send_date IS NULL OR send_date <= ?)",
				campaign.ID, "scheduled", now).
			Preload("Recipient").
			Find(&results).Error; err != nil {
			return err
		}
		if len(results) == 0 {
			return nil
		}
		ids := make([]int64, len(results))
		for i := range results {
			ids[i] = results[i].ID
		}
		return tx.Model(&model.Result{}).Where("id IN ?", ids).Update("status", model.ResultStatusSending).Error
	})
	if err != nil {
		log.Printf("campaign %d: claim results error: %v", campaign.ID, err)
		return
	}

	if len(results) == 0 {
		// Check if all results are sent/error — mark campaign complete
		var pending int64
		database.Model(&model.Result{}).
			Where("campaign_id = ? AND status IN ?", campaign.ID, []string{"scheduled", model.ResultStatusSending}).
			Count(&pending)
		if pending == 0 {
			now := time.Now().UTC()
			campaign.Status = model.CampaignStatusCompleted
			campaign.CompletedAt = &now
			campaignRepo.Update(campaign)
			log.Printf("campaign %d: completed", campaign.ID)
			// Auto-send report to tenant admin
			go sendCompletionReport(database, cfg, resultRepo, campaign)
		}
		return
	}

	log.Printf("campaign %d: sending %d emails", campaign.ID, len(results))

	// Create rate limiter based on provider type
	rl := getRateLimiter(smtp.MailerType)

	for i := range results {
		r := &results[i]
		if r.Recipient == nil {
			continue
		}

		// Check domain rate limit
		recipientDomain := extractDomain(r.Recipient.Email)
		if !rl.Wait(recipientDomain) {
			// Domain hourly limit reached — revert to scheduled for next cycle
			database.Model(r).Update("status", "scheduled")
			log.Printf("campaign %d: domain %s hourly limit reached, deferring %s",
				campaign.ID, recipientDomain, r.Recipient.Email)
			continue
		}

		// Render email with tracking URLs
		rid := r.RID
		if rid == "" {
			rid = uuid.New().String()
			r.RID = rid
			database.Model(r).Update("rid", rid)
		}
		trackBase := cfg.TrackerBaseURL
		htmlBody := renderTemplate(tmpl.HTMLBody, r.Recipient, trackBase, rid)

		// Compliance headers
		reportURL := fmt.Sprintf("%s/t/r/%s", trackBase, rid)
		headers := map[string]string{
			"List-Unsubscribe":      "<" + reportURL + ">",
			"List-Unsubscribe-Post": "List-Unsubscribe=One-Click",
			"X-Mailer":              "PhishGuard/1.0",
			"Precedence":            "bulk",
			"Message-ID":            fmt.Sprintf("<%s@%s>", rid, extractDomain(smtp.FromAddress)),
		}

		msg := &mailer.Message{
			From:     smtp.FromAddress,
			FromName: smtp.FromName,
			To:       r.Recipient.Email,
			Subject:  tmpl.Subject,
			HTMLBody: htmlBody,
			TextBody: tmpl.TextBody,
			Headers:  headers,
		}

		if err := m.Send(context.Background(), msg); err != nil {
			log.Printf("campaign %d: failed to send to %s: %v", campaign.ID, r.Recipient.Email, err)
			if isTransientError(err) {
				// Revert to scheduled for retry next cycle
				database.Model(r).Update("status", "scheduled")
				continue
			}
			// Permanent error — mark as failed
			database.Model(r).Select("status", "error_detail").Updates(model.Result{
				Status: model.EventError, ErrorDetail: err.Error(),
			})
		} else {
			sentAt := time.Now().UTC()
			database.Model(r).Select("status", "sent_at").Updates(model.Result{
				Status: model.EventSent, SentAt: &sentAt,
			})
		}
	}
}

func isTransientError(err error) bool {
	msg := strings.ToLower(err.Error())
	transientPatterns := []string{
		// SMTP transient codes
		"421", "450", "451", "452",
		// Connection issues
		"connection reset", "timeout", "eof",
		// AWS SES specific
		"throttlingexception", "throttling", "rate exceeded",
		"toomanyrequestsexception", "maximum sending rate",
		// Mailgun specific
		"429", "rate limit", "too many requests",
		// Generic
		"temporary", "try again",
	}
	for _, p := range transientPatterns {
		if strings.Contains(msg, p) {
			return true
		}
	}
	return false
}

func extractDomain(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return "phishguard.local"
}

func renderTemplate(html string, recipient *model.Recipient, trackBase, rid string) string {
	// Simple placeholder replacement
	replacer := map[string]string{
		"{{.FirstName}}":  recipient.FirstName,
		"{{.LastName}}":   recipient.LastName,
		"{{.Email}}":      recipient.Email,
		"{{.Department}}": recipient.Department,
		"{{.Position}}":   recipient.Position,
	}
	result := html
	for k, v := range replacer {
		result = replaceAll(result, k, v)
	}

	// Append tracking pixel
	pixel := fmt.Sprintf(`<img src="%s/t/o/%s" width="1" height="1" style="display:none" />`, trackBase, rid)
	result += pixel

	// Replace links with tracking URLs
	trackClick := fmt.Sprintf("%s/t/c/%s", trackBase, rid)
	result = replaceAll(result, "{{.TrackURL}}", trackClick)

	// Replace download URL
	downloadURL := fmt.Sprintf("%s/t/d/%s/attachment", trackBase, rid)
	result = replaceAll(result, "{{.DownloadURL}}", downloadURL)

	// Replace report URL
	reportURL := fmt.Sprintf("%s/t/r/%s", trackBase, rid)
	result = replaceAll(result, "{{.ReportURL}}", reportURL)

	// Append report link at bottom
	reportLink := fmt.Sprintf(`<div style="margin-top:32px;padding-top:12px;border-top:1px solid #eee;font-size:11px;color:#999;text-align:center;">覺得這封信可疑？<a href="%s" style="color:#999;">點此舉報</a></div>`, reportURL)
	result += reportLink

	return result
}

func replaceAll(s, old, new string) string {
	for {
		i := indexOf(s, old)
		if i < 0 {
			return s
		}
		s = s[:i] + new + s[i+len(old):]
	}
}

func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func buildMailer(smtp *model.SMTPProfile, encryptKey string) (mailer.Mailer, error) {
	cfg := map[string]string{
		"from_address": smtp.FromAddress,
		"from_name":    smtp.FromName,
	}
	switch smtp.MailerType {
	case "smtp":
		cfg["host"] = smtp.Host
		if smtp.Port != nil {
			cfg["port"] = fmt.Sprintf("%d", *smtp.Port)
		}
		cfg["username"] = smtp.Username
		pw, err := crypto.Decrypt(encryptKey, smtp.PasswordEnc)
		if err != nil {
			return nil, fmt.Errorf("decrypt password: %w", err)
		}
		cfg["password"] = pw
		if smtp.TLSRequired {
			cfg["tls"] = "true"
		}
	case "mailgun":
		cfg["domain"] = smtp.MailgunDomain
		apiKey, err := crypto.Decrypt(encryptKey, smtp.MailgunAPIKey)
		if err != nil {
			return nil, fmt.Errorf("decrypt mailgun api key: %w", err)
		}
		cfg["api_key"] = apiKey
	case "ses":
		cfg["region"] = smtp.SESRegion
		ak, err := crypto.Decrypt(encryptKey, smtp.SESAccessKey)
		if err != nil {
			return nil, fmt.Errorf("decrypt ses access key: %w", err)
		}
		cfg["access_key"] = ak
		sk, err := crypto.Decrypt(encryptKey, smtp.SESSecretKey)
		if err != nil {
			return nil, fmt.Errorf("decrypt ses secret key: %w", err)
		}
		cfg["secret_key"] = sk
	}
	return mailer.NewMailer(smtp.MailerType, cfg)
}

// Ensure asynq is used (for future queue-based processing)
func runScheduler(database *gorm.DB, campaignSvc *service.CampaignService) {
	if err := service.RunScheduledTests(database, campaignSvc); err != nil {
		log.Printf("scheduler error: %v", err)
	}
}

// sendCompletionReport sends report email when campaign just completed
func sendCompletionReport(database *gorm.DB, cfg *config.Config, resultRepo *repo.ResultRepo, campaign *model.Campaign) {
	var smtp model.SMTPProfile
	if err := database.First(&smtp, campaign.SMTPProfileID).Error; err != nil {
		return
	}
	m, err := buildMailer(&smtp, cfg.EncryptKey)
	if err != nil {
		return
	}
	if err := service.SendCampaignReport(database, resultRepo, campaign, m, smtp.FromAddress); err != nil {
		log.Printf("campaign %d: failed to send report: %v", campaign.ID, err)
	} else {
		log.Printf("campaign %d: report sent to tenant admin", campaign.ID)
	}
}

var _ = asynq.NewClient
