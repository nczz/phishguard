package service

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/phishguard/phishguard/internal/model"
	"gorm.io/gorm"
)

func RunScheduledTests(db *gorm.DB, campaignSvc *CampaignService) error {
	var configs []model.AutoTestConfig
	if err := db.Where("is_enabled = ? AND next_run_at <= ?", true, time.Now()).Find(&configs).Error; err != nil {
		return fmt.Errorf("query auto_test_configs: %w", err)
	}

	for _, cfg := range configs {
		if err := runOneAutoTest(db, campaignSvc, &cfg); err != nil {
			log.Printf("[auto-test] tenant=%d error: %v", cfg.TenantID, err)
			continue
		}
	}
	return nil
}

func runOneAutoTest(db *gorm.DB, campaignSvc *CampaignService, cfg *model.AutoTestConfig) error {
	// Tenant
	var tenant model.Tenant
	if err := db.First(&tenant, cfg.TenantID).Error; err != nil {
		return fmt.Errorf("tenant %d not found: %w", cfg.TenantID, err)
	}

	// Random active scenario
	var scenario model.Scenario
	if err := db.Where("(tenant_id = ? OR tenant_id IS NULL) AND is_active = ?", cfg.TenantID, true).
		Order("RAND()").First(&scenario).Error; err != nil {
		return fmt.Errorf("no active scenario: %w", err)
	}

	// All recipient groups
	var groups []model.RecipientGroup
	if err := db.Where("tenant_id = ?", cfg.TenantID).Find(&groups).Error; err != nil {
		return fmt.Errorf("find groups: %w", err)
	}
	if len(groups) == 0 {
		return fmt.Errorf("no recipient groups for tenant %d", cfg.TenantID)
	}

	// First SMTP profile
	var smtp model.SMTPProfile
	if err := db.Where("tenant_id = ?", cfg.TenantID).First(&smtp).Error; err != nil {
		return fmt.Errorf("no SMTP profile: %w", err)
	}

	groupIDs := make([]int64, len(groups))
	for i, g := range groups {
		groupIDs[i] = g.ID
	}

	phishURL := os.Getenv("PHISH_URL")
	if phishURL == "" {
		phishURL = "https://phish.example.com"
	}

	scenarioID := scenario.ID
	campaign, err := campaignSvc.CreateCampaign(cfg.TenantID, &CreateCampaignRequest{
		Name:          fmt.Sprintf("Auto Test - %s", time.Now().Format("2006-01-02")),
		ScenarioID:    &scenarioID,
		SMTPProfileID: smtp.ID,
		GroupIDs:      groupIDs,
		PhishURL:      phishURL,
		SelectionMode: cfg.TargetMode,
		SamplePercent: cfg.SamplePercent,
	})
	if err != nil {
		return fmt.Errorf("create campaign: %w", err)
	}

	if err := campaignSvc.LaunchCampaign(cfg.TenantID, campaign.ID); err != nil {
		return fmt.Errorf("launch campaign: %w", err)
	}

	// Calculate next run
	now := time.Now()
	switch cfg.Frequency {
	case "monthly":
		now = now.AddDate(0, 1, 0)
	case "quarterly":
		now = now.AddDate(0, 3, 0)
	case "biannual":
		now = now.AddDate(0, 6, 0)
	default:
		now = now.AddDate(0, 1, 0)
	}
	cfg.NextRunAt = &now
	if err := db.Model(cfg).Update("next_run_at", cfg.NextRunAt).Error; err != nil {
		return fmt.Errorf("update next_run_at: %w", err)
	}

	log.Printf("[auto-test] tenant=%d campaign=%d launched, next=%s", cfg.TenantID, campaign.ID, now.Format("2006-01-02"))
	return nil
}
