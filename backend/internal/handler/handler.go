package handler

import (
	"github.com/phishguard/phishguard/internal/repo"
	"github.com/phishguard/phishguard/internal/service"
)

type Handler struct {
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
