package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phishguard/phishguard/internal/model"
)

func (h *Handler) CreateTenant(c *gin.Context) {
	var req struct {
		Name          string `json:"name"`
		Slug          string `json:"slug"`
		Plan          string `json:"plan"`
		AdminEmail    string `json:"admin_email"`
		AdminPassword string `json:"admin_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tenant, err := h.TenantService.CreateWithAdmin(req.Name, req.Slug, req.Plan, req.AdminEmail, req.AdminPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tenant)
}

func (h *Handler) ListTenants(c *gin.Context) {
	tenants, err := h.TenantRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tenants)
}

func (h *Handler) GetTenant(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	tenant, err := h.TenantRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}
	c.JSON(http.StatusOK, tenant)
}

func (h *Handler) UpdateTenant(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	tenant, err := h.TenantRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}
	var req struct {
		Name          *string `json:"name"`
		Plan          *string `json:"plan"`
		IsActive      *bool   `json:"is_active"`
		MaxRecipients *int    `json:"max_recipients"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name != nil {
		tenant.Name = *req.Name
	}
	if req.Plan != nil {
		tenant.Plan = *req.Plan
	}
	if req.IsActive != nil {
		tenant.IsActive = *req.IsActive
	}
	if req.MaxRecipients != nil {
		tenant.MaxRecipients = *req.MaxRecipients
	}
	if err := h.TenantRepo.Update(tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tenant)
}

func (h *Handler) AdminDashboard(c *gin.Context) {
	stats, err := h.TenantService.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) CreatePlatformScenario(c *gin.Context) {
	var scenario model.Scenario
	if err := c.ShouldBindJSON(&scenario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	scenario.TenantID = nil
	if err := h.ScenarioRepo.Create(&scenario); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, scenario)
}
