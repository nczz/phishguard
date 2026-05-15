package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/middleware"
	"github.com/nczz/phishguard/internal/model"
	"github.com/nczz/phishguard/internal/service"
)

// ListPlans returns all plan definitions
func (h *Handler) ListPlans(c *gin.Context) {
	c.JSON(http.StatusOK, service.PlanDefaults)
}

// AdminGetPlan returns a tenant's effective plan config + current usage
func (h *Handler) AdminGetPlan(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var tenant model.Tenant
	if err := h.DB.First(&tenant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}
	limits := service.GetEffectiveLimits(&tenant)
	recipientCount, _ := h.RecipientRepo.CountByTenant(id)
	campaignCount, _ := h.CampaignRepo.CountByTenantThisYear(id)
	emailsThisMonth, _ := h.ResultRepo.CountSentThisMonth(id)

	c.JSON(http.StatusOK, gin.H{
		"plan":          tenant.Plan,
		"limits":        limits,
		"plan_defaults": service.GetPlanConfig(tenant.Plan),
		"overrides":     gin.H{"max_recipients": tenant.MaxRecipients, "max_campaigns_per_year": tenant.MaxCampaignsPerYear},
		"usage":         gin.H{"recipients": recipientCount, "campaigns_this_year": campaignCount, "emails_this_month": emailsThisMonth},
	})
}

// AdminUpdatePlan updates a tenant's plan and/or overrides
func (h *Handler) AdminUpdatePlan(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if _, err := h.TenantRepo.FindByID(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}
	var req struct {
		Plan                string `json:"plan"`
		MaxRecipients       *int   `json:"max_recipients"`
		MaxCampaignsPerYear *int   `json:"max_campaigns_per_year"`
		MaxEmailsPerMonth   *int   `json:"max_emails_per_month"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{}
	if req.Plan != "" {
		if _, ok := service.PlanDefaults[req.Plan]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan"})
			return
		}
		updates["plan"] = req.Plan
	}
	if req.MaxRecipients != nil {
		if *req.MaxRecipients < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "max_recipients must be >= 0"})
			return
		}
		updates["max_recipients"] = *req.MaxRecipients
	}
	if req.MaxCampaignsPerYear != nil {
		if *req.MaxCampaignsPerYear < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "max_campaigns_per_year must be >= 0"})
			return
		}
		updates["max_campaigns_per_year"] = *req.MaxCampaignsPerYear
	}
	if req.MaxEmailsPerMonth != nil {
		if *req.MaxEmailsPerMonth < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "max_emails_per_month must be >= 0"})
			return
		}
		updates["max_emails_per_month"] = *req.MaxEmailsPerMonth
	}
	h.DB.Model(&model.Tenant{}).Where("id = ?", id).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// GetMyPlan returns the current tenant's plan info + usage (for tenant users)
func (h *Handler) GetMyPlan(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	var tenant model.Tenant
	if err := h.DB.First(&tenant, tid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}
	limits := service.GetEffectiveLimits(&tenant)
	recipientCount, _ := h.RecipientRepo.CountByTenant(tid)
	campaignCount, _ := h.CampaignRepo.CountByTenantThisYear(tid)
	emailsThisMonth, _ := h.ResultRepo.CountSentThisMonth(tid)

	c.JSON(http.StatusOK, gin.H{
		"plan":   tenant.Plan,
		"limits": limits,
		"usage":  gin.H{"recipients": recipientCount, "campaigns_this_year": campaignCount, "emails_this_month": emailsThisMonth},
	})
}
