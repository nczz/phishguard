package main

import (
	"log"

	"github.com/nczz/phishguard/config"
	"github.com/nczz/phishguard/internal/db"
	"github.com/nczz/phishguard/internal/handler"
	"github.com/nczz/phishguard/internal/repo"
	"github.com/nczz/phishguard/internal/service"
)

func main() {
	cfg := config.Load()
	database := db.Init(cfg.DBDSN)

	tenantRepo := &repo.TenantRepo{DB: database}
	userRepo := &repo.UserRepo{DB: database}
	templateRepo := &repo.TemplateRepo{DB: database}
	pageRepo := &repo.PageRepo{DB: database}
	scenarioRepo := &repo.ScenarioRepo{DB: database}
	recipientRepo := &repo.RecipientRepo{DB: database}
	smtpRepo := &repo.SMTPRepo{DB: database}
	campaignRepo := &repo.CampaignRepo{DB: database}
	resultRepo := &repo.ResultRepo{DB: database}
	auditRepo := &repo.AuditRepo{DB: database}

	authSvc := &service.AuthService{UserRepo: userRepo, JWTSecret: cfg.JWTSecret}
	tenantSvc := &service.TenantService{TenantRepo: tenantRepo, UserRepo: userRepo, DB: database}
	campaignSvc := &service.CampaignService{
		CampaignRepo:  campaignRepo,
		ResultRepo:    resultRepo,
		RecipientRepo: recipientRepo,
		ScenarioRepo:  scenarioRepo,
	}
	reportSvc := &service.ReportService{ResultRepo: resultRepo}

	// Create initial admin if not exists
	if err := authSvc.CreateInitialAdmin(cfg.AdminEmail, cfg.AdminPassword); err != nil {
		log.Printf("initial admin: %v", err)
	}

	h := &handler.Handler{
		DB:              database,
		JWTSecret:       cfg.JWTSecret,
		AuthService:     authSvc,
		TenantService:   tenantSvc,
		CampaignService: campaignSvc,
		ReportService:   reportSvc,
		TenantRepo:      tenantRepo,
		UserRepo:        userRepo,
		TemplateRepo:    templateRepo,
		PageRepo:        pageRepo,
		ScenarioRepo:    scenarioRepo,
		RecipientRepo:   recipientRepo,
		SMTPRepo:        smtpRepo,
		CampaignRepo:    campaignRepo,
		ResultRepo:      resultRepo,
		AuditRepo:       auditRepo,
	}

	auditLogger := &repo.DBAuditLogger{Repo: auditRepo}
	r := handler.SetupRouter(h, cfg.JWTSecret, auditLogger)
	log.Printf("API server starting on %s", cfg.APIAddr)
	if err := r.Run(cfg.APIAddr); err != nil {
		log.Fatalf("failed to start API server: %v", err)
	}
}
