package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/middleware"
	"github.com/nczz/phishguard/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// --- Tenant Stats ---

type TenantStats struct {
	TenantID       int64  `json:"tenant_id"`
	TenantName     string `json:"tenant_name"`
	Slug           string `json:"slug"`
	Plan           string `json:"plan"`
	IsActive       bool   `json:"is_active"`
	RecipientCount int64  `json:"recipient_count"`
	CampaignCount  int64  `json:"campaign_count"`
	EmailsSent     int64  `json:"emails_sent"`
}

func (h *Handler) AdminDashboardFull(c *gin.Context) {
	tenants, err := h.TenantRepo.FindAll()
	if err != nil {
		serverError(c, err)
		return
	}

	stats := make([]TenantStats, 0, len(tenants))
	var totalRecipients, totalCampaigns, totalEmails int64

	for _, t := range tenants {
		var rCount, cCount, eCount int64
		h.DB.Model(&model.Recipient{}).Where("tenant_id = ?", t.ID).Count(&rCount)
		h.DB.Model(&model.Campaign{}).Where("tenant_id = ?", t.ID).Count(&cCount)
		h.DB.Model(&model.Result{}).Where("tenant_id = ? AND sent_at IS NOT NULL", t.ID).Count(&eCount)

		stats = append(stats, TenantStats{
			TenantID: t.ID, TenantName: t.Name, Slug: t.Slug,
			Plan: t.Plan, IsActive: t.IsActive,
			RecipientCount: rCount, CampaignCount: cCount, EmailsSent: eCount,
		})
		totalRecipients += rCount
		totalCampaigns += cCount
		totalEmails += eCount
	}

	// Alerts
	alerts := []gin.H{}
	for _, s := range stats {
		if s.IsActive && s.CampaignCount == 0 {
			alerts = append(alerts, gin.H{"type": "warning", "message": s.TenantName + " 尚未建立任何測試活動"})
		}
		if !s.IsActive {
			alerts = append(alerts, gin.H{"type": "error", "message": s.TenantName + " 已停用"})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_tenants":    len(tenants),
		"active_tenants":   countActive(tenants),
		"total_recipients": totalRecipients,
		"total_campaigns":  totalCampaigns,
		"total_emails":     totalEmails,
		"tenants":          stats,
		"alerts":           alerts,
	})
}

func countActive(tenants []model.Tenant) int {
	n := 0
	for _, t := range tenants {
		if t.IsActive {
			n++
		}
	}
	return n
}

// --- Tenant Campaigns (cross-tenant view) ---

func (h *Handler) AdminTenantCampaigns(c *gin.Context) {
	tid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var campaigns []model.Campaign
	h.DB.Where("tenant_id = ?", tid).Order("created_at DESC").Find(&campaigns)
	c.JSON(http.StatusOK, campaigns)
}

// --- Toggle Tenant Active ---

func (h *Handler) ToggleTenant(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.DB.Model(&model.Tenant{}).Where("id = ?", id).Update("is_active", req.IsActive)
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// --- User Management ---

func (h *Handler) AdminListUsers(c *gin.Context) {
	tid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	users, err := h.UserRepo.FindAllByTenant(tid)
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) AdminCreateUser(c *gin.Context) {
	tid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		Email    string `json:"email" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	u := &model.User{TenantID: &tid, Email: req.Email, Name: req.Name, PasswordHash: string(hash), Role: req.Role, IsActive: true}
	if err := h.UserRepo.Create(u); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, u)
}

func (h *Handler) AdminDeleteUser(c *gin.Context) {
	tid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	uid, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	user, err := h.UserRepo.FindByID(uid)
	if err != nil {
		notFound(c, "user not found")
		return
	}
	if user.TenantID == nil || *user.TenantID != tid {
		forbidden(c, "user does not belong to this tenant")
		return
	}
	if err := h.UserRepo.Delete(uid); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) AdminUpdateUser(c *gin.Context) {
	uid, _ := strconv.ParseInt(c.Param("uid"), 10, 64)
	var req struct {
		Name     string `json:"name"`
		Role     string `json:"role"`
		IsActive *bool  `json:"is_active"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.UserRepo.FindByID(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		user.PasswordHash = string(hash)
	}
	if err := h.UserRepo.Update(user); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// --- Cross-tenant Audit Logs ---

func (h *Handler) AdminAuditLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	logs, total, err := h.AuditRepo.FindAll(limit, offset)
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs, "total": total})
}

// --- Impersonate (generate token for tenant) ---

func (h *Handler) AdminImpersonate(c *gin.Context) {
	tid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	admin := middleware.GetClaims(c)

	token, err := middleware.GenerateToken(h.JWTSecret, admin.UserID, &tid, "tenant_admin", admin.Email)
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "tenant_id": tid})
}
