package service

import (
	"context"
	"fmt"
	"time"

	"github.com/nczz/phishguard/internal/mailer"
	"github.com/nczz/phishguard/internal/model"
	"github.com/nczz/phishguard/internal/repo"
	"gorm.io/gorm"
)

// SendCampaignReport sends a campaign completion report email to the tenant admin.
func SendCampaignReport(db *gorm.DB, resultRepo *repo.ResultRepo, campaign *model.Campaign, mailerSvc mailer.Mailer, fromAddr string) error {
	stats, err := resultRepo.GetFunnelStats(campaign.TenantID, campaign.ID)
	if err != nil {
		return fmt.Errorf("get funnel stats: %w", err)
	}

	html := buildReportHTML(campaign, stats)

	// Find tenant admin email
	var admin model.User
	if err := db.Where("tenant_id = ? AND role = ?", campaign.TenantID, "tenant_admin").First(&admin).Error; err != nil {
		return fmt.Errorf("find tenant admin: %w", err)
	}

	msg := &mailer.Message{
		From:     fromAddr,
		FromName: "PhishGuard",
		To:       admin.Email,
		Subject:  "[PhishGuard] Campaign Report: " + campaign.Name,
		HTMLBody: html,
	}
	if err := mailerSvc.Send(context.Background(), msg); err != nil {
		return fmt.Errorf("send report email: %w", err)
	}
	return nil
}

func buildReportHTML(c *model.Campaign, s *repo.FunnelStats) string {
	completed := time.Now().Format("2006-01-02 15:04")
	if c.CompletedAt != nil {
		completed = c.CompletedAt.Format("2006-01-02 15:04")
	}

	pct := func(n int64) string {
		if s.Total == 0 {
			return "0.0%"
		}
		return fmt.Sprintf("%.1f%%", float64(n)/float64(s.Total)*100)
	}

	return fmt.Sprintf(`<html><body style="font-family:sans-serif;color:#333">
<h2>Campaign Report: %s</h2>
<p>Completed: %s</p>
<table border="1" cellpadding="8" cellspacing="0" style="border-collapse:collapse">
<tr style="background:#f5f5f5"><th>Stage</th><th>Count</th><th>Rate</th></tr>
<tr><td>Sent</td><td>%d</td><td>%s</td></tr>
<tr><td>Opened</td><td>%d</td><td>%s</td></tr>
<tr><td>Clicked</td><td>%d</td><td>%s</td></tr>
<tr><td>Submitted</td><td>%d</td><td>%s</td></tr>
<tr><td>Reported</td><td>%d</td><td>%s</td></tr>
</table>
<p style="margin-top:16px;color:#666">Login to PhishGuard for full report with department breakdown and recipient details</p>
</body></html>`,
		c.Name, completed,
		s.Sent, pct(s.Sent),
		s.Opened, pct(s.Opened),
		s.Clicked, pct(s.Clicked),
		s.Submitted, pct(s.Submitted),
		s.Reported, pct(s.Reported),
	)
}
