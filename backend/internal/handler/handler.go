package handler

import (
	"github.com/nczz/phishguard/internal/repo"
	"github.com/nczz/phishguard/internal/service"
	"gorm.io/gorm"
)

type Handler struct {
	DB              *gorm.DB
	JWTSecret       string
	AuthService     *service.AuthService
	TenantService   *service.TenantService
	CampaignService *service.CampaignService
	ReportService   *service.ReportService
	TenantRepo      *repo.TenantRepo
	UserRepo        *repo.UserRepo
	TemplateRepo    *repo.TemplateRepo
	PageRepo        *repo.PageRepo
	ScenarioRepo    *repo.ScenarioRepo
	RecipientRepo   *repo.RecipientRepo
	SMTPRepo        *repo.SMTPRepo
	CampaignRepo    *repo.CampaignRepo
	ResultRepo      *repo.ResultRepo
	AuditRepo       *repo.AuditRepo
}
