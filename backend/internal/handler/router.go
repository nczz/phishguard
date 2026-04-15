package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/phishguard/phishguard/internal/middleware"
)

func SetupRouter(h *Handler, jwtSecret string) *gin.Engine {
	r := gin.Default()

	auth := r.Group("/api/auth")
	{
		auth.POST("/login", h.Login)
		auth.GET("/me", middleware.AuthMiddleware(jwtSecret), h.Me)
	}

	admin := r.Group("/api/admin",
		middleware.AuthMiddleware(jwtSecret),
		middleware.RoleRequired("platform_admin"),
		middleware.TenantMiddleware(),
	)
	{
		admin.POST("/tenants", h.CreateTenant)
		admin.GET("/tenants", h.ListTenants)
		admin.GET("/tenants/:id", h.GetTenant)
		admin.PUT("/tenants/:id", h.UpdateTenant)
		admin.GET("/dashboard", h.AdminDashboard)
		admin.POST("/scenarios", h.CreatePlatformScenario)
	}

	api := r.Group("/api",
		middleware.AuthMiddleware(jwtSecret),
		middleware.TenantMiddleware(),
		middleware.RequireTenant(),
	)
	{
		api.GET("/scenarios", h.ListScenarios)
		api.POST("/scenarios", h.CreateScenario)

		api.GET("/templates", h.ListTemplates)
		api.POST("/templates", h.CreateTemplate)
		api.PUT("/templates/:id", h.UpdateTemplate)
		api.DELETE("/templates/:id", h.DeleteTemplate)

		api.GET("/pages", h.ListPages)
		api.POST("/pages", h.CreatePage)
		api.PUT("/pages/:id", h.UpdatePage)
		api.DELETE("/pages/:id", h.DeletePage)

		api.GET("/recipient-groups", h.ListRecipientGroups)
		api.POST("/recipient-groups", h.CreateRecipientGroup)
		api.POST("/recipient-groups/import", h.ImportRecipients)

		api.GET("/smtp-profiles", h.ListSMTPProfiles)
		api.POST("/smtp-profiles", h.CreateSMTPProfile)
		api.POST("/smtp-profiles/:id/test", h.TestSMTPProfile)

		api.POST("/campaigns", h.CreateCampaign)
		api.GET("/campaigns", h.ListCampaigns)
		api.GET("/campaigns/:id", h.GetCampaign)
		api.POST("/campaigns/:id/launch", h.LaunchCampaign)
		api.DELETE("/campaigns/:id", h.DeleteCampaign)

		api.GET("/campaigns/:id/report", h.GetCampaignReport)
		api.GET("/campaigns/:id/report/pdf", h.ExportCampaignPDF)

		api.GET("/reports/overview", h.GetOverviewReport)
		api.GET("/reports/department", h.GetDepartmentReport)

		api.GET("/audit-logs", h.ListAuditLogs)
	}

	return r
}
