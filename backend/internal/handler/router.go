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
		admin.GET("/dashboard", h.AdminDashboardFull)
		admin.POST("/scenarios", h.CreatePlatformScenario)
		admin.GET("/tenants/:id/campaigns", h.AdminTenantCampaigns)
		admin.PATCH("/tenants/:id/toggle", h.ToggleTenant)
		admin.GET("/tenants/:id/users", h.AdminListUsers)
		admin.POST("/tenants/:id/users", h.AdminCreateUser)
		admin.DELETE("/tenants/:id/users/:uid", h.AdminDeleteUser)
		admin.PUT("/tenants/:id/users/:uid", h.AdminUpdateUser)
		admin.POST("/tenants/:id/impersonate", h.AdminImpersonate)
		admin.GET("/tenants/:id/plan", h.AdminGetPlan)
		admin.PUT("/tenants/:id/plan", h.AdminUpdatePlan)
		admin.GET("/plans", h.ListPlans)
		admin.GET("/audit-logs", h.AdminAuditLogs)
	}

	api := r.Group("/api",
		middleware.AuthMiddleware(jwtSecret),
		middleware.TenantMiddleware(),
		middleware.RequireTenant(),
	)
	{
		api.GET("/scenarios", h.ListScenarios)
		api.POST("/scenarios", h.CreateScenario)
		api.PUT("/scenarios/:id", h.UpdateScenario)
		api.DELETE("/scenarios/:id", h.DeleteScenario)

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
		api.PUT("/recipients/:id", h.UpdateRecipient)
		api.DELETE("/recipients/:id", h.DeleteRecipient)
		api.POST("/recipients/batch-delete", h.BatchDeleteRecipients)

		api.GET("/smtp-profiles", h.ListSMTPProfiles)
		api.POST("/smtp-profiles", h.CreateSMTPProfile)
		api.POST("/smtp-profiles/:id/test", h.TestSMTPProfile)
		api.POST("/smtp-profiles/check-compliance", h.CheckMailCompliance)

		api.POST("/campaigns", h.CreateCampaign)
		api.GET("/campaigns", h.ListCampaigns)
		api.GET("/campaigns/:id", h.GetCampaign)
		api.POST("/campaigns/:id/launch", h.LaunchCampaign)
		api.DELETE("/campaigns/:id", h.DeleteCampaign)

		api.GET("/campaigns/:id/report", h.GetCampaignReport)
		api.GET("/campaigns/:id/report/pdf", h.ExportCampaignPDFReal)
		api.GET("/campaigns/:id/recipients", h.CampaignRecipients)
		api.GET("/campaigns/:id/export/csv", h.ExportCampaignCSV)

		api.GET("/reports/overview", h.GetOverviewReport)
		api.GET("/reports/department", h.GetDepartmentReport)
		api.GET("/reports/dashboard-stats", h.TenantDashboardStats)
		api.GET("/reports/offenders", h.RepeatOffenders)
		api.GET("/reports/trend", h.TrendAnalysis)

		api.POST("/campaigns/:id/send-report", h.SendCampaignReportEmail)

		api.GET("/audit-logs", h.ListAuditLogs)
		api.POST("/seed-sample-data", h.SeedSampleData)
		api.GET("/my-plan", h.GetMyPlan)
		api.GET("/auto-test", h.GetAutoTestConfig)
		api.PUT("/auto-test", h.SaveAutoTestConfig)
	}

	return r
}
