package db

import (
	"log"

	"github.com/phishguard/phishguard/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(
		&model.Tenant{},
		&model.User{},
		&model.EmailTemplate{},
		&model.LandingPage{},
		&model.Scenario{},
		&model.RecipientGroup{},
		&model.Recipient{},
		&model.SMTPProfile{},
		&model.Campaign{},
		&model.CampaignGroup{},
		&model.Result{},
		&model.Event{},
		&model.AutoTestConfig{},
		&model.Subscription{},
		&model.UsageRecord{},
		&model.AuditLog{},
	); err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}

	return db
}
