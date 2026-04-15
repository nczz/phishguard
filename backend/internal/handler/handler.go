package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB        *gorm.DB
	JWTSecret string
}

func NewHandler(db *gorm.DB, jwtSecret string) *Handler {
	return &Handler{DB: db, JWTSecret: jwtSecret}
}

var stub = func(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// Auth
func (h *Handler) Login(c *gin.Context)    { stub(c) }
func (h *Handler) Me(c *gin.Context)       { stub(c) }

// Platform Admin
func (h *Handler) CreateTenant(c *gin.Context)           { stub(c) }
func (h *Handler) ListTenants(c *gin.Context)            { stub(c) }
func (h *Handler) GetTenant(c *gin.Context)              { stub(c) }
func (h *Handler) UpdateTenant(c *gin.Context)           { stub(c) }
func (h *Handler) AdminDashboard(c *gin.Context)         { stub(c) }
func (h *Handler) CreatePlatformScenario(c *gin.Context) { stub(c) }

// Scenarios
func (h *Handler) ListScenarios(c *gin.Context)  { stub(c) }
func (h *Handler) CreateScenario(c *gin.Context) { stub(c) }

// Templates
func (h *Handler) ListTemplates(c *gin.Context)  { stub(c) }
func (h *Handler) CreateTemplate(c *gin.Context) { stub(c) }
func (h *Handler) UpdateTemplate(c *gin.Context) { stub(c) }
func (h *Handler) DeleteTemplate(c *gin.Context) { stub(c) }

// Pages
func (h *Handler) ListPages(c *gin.Context)  { stub(c) }
func (h *Handler) CreatePage(c *gin.Context) { stub(c) }
func (h *Handler) UpdatePage(c *gin.Context) { stub(c) }
func (h *Handler) DeletePage(c *gin.Context) { stub(c) }

// Recipients
func (h *Handler) ListRecipientGroups(c *gin.Context)  { stub(c) }
func (h *Handler) CreateRecipientGroup(c *gin.Context) { stub(c) }
func (h *Handler) ImportRecipients(c *gin.Context)     { stub(c) }

// SMTP
func (h *Handler) ListSMTPProfiles(c *gin.Context)  { stub(c) }
func (h *Handler) CreateSMTPProfile(c *gin.Context) { stub(c) }
func (h *Handler) TestSMTPProfile(c *gin.Context)   { stub(c) }

// Campaigns
func (h *Handler) CreateCampaign(c *gin.Context) { stub(c) }
func (h *Handler) ListCampaigns(c *gin.Context)  { stub(c) }
func (h *Handler) GetCampaign(c *gin.Context)    { stub(c) }
func (h *Handler) LaunchCampaign(c *gin.Context) { stub(c) }
func (h *Handler) DeleteCampaign(c *gin.Context) { stub(c) }

// Reports
func (h *Handler) GetCampaignReport(c *gin.Context)   { stub(c) }
func (h *Handler) ExportCampaignPDF(c *gin.Context)    { stub(c) }
func (h *Handler) GetOverviewReport(c *gin.Context)    { stub(c) }
func (h *Handler) GetDepartmentReport(c *gin.Context)  { stub(c) }

// Audit
func (h *Handler) ListAuditLogs(c *gin.Context) { stub(c) }
