package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/middleware"
	"github.com/nczz/phishguard/internal/model"
	"github.com/nczz/phishguard/internal/service"
)

func (h *Handler) GetAutoTestConfig(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	var cfg model.AutoTestConfig
	if err := h.DB.Where("tenant_id = ?", tid).First(&cfg).Error; err != nil {
		// Return default
		c.JSON(http.StatusOK, model.AutoTestConfig{TenantID: tid, Frequency: "quarterly", TargetMode: "random", SamplePercent: 30})
		return
	}
	c.JSON(http.StatusOK, cfg)
}

func (h *Handler) SaveAutoTestConfig(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	var req struct {
		IsEnabled     bool   `json:"is_enabled"`
		Frequency     string `json:"frequency"`
		TargetMode    string `json:"target_mode"`
		SamplePercent int    `json:"sample_percent"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Frequency == "" {
		req.Frequency = "quarterly"
	}
	if req.Frequency != "monthly" && req.Frequency != "quarterly" && req.Frequency != "biannual" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid frequency"})
		return
	}
	if req.TargetMode == "" {
		req.TargetMode = "random"
	}
	if req.TargetMode != "all" && req.TargetMode != "random" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target_mode"})
		return
	}
	if req.SamplePercent < 1 || req.SamplePercent > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sample_percent must be 1-100"})
		return
	}
	if req.IsEnabled {
		var tenant model.Tenant
		if err := h.DB.First(&tenant, tid).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
			return
		}
		if !service.GetEffectiveLimits(&tenant).AutoTest {
			c.JSON(http.StatusForbidden, gin.H{"error": "目前方案不支援自動測試，請升級方案"})
			return
		}
	}

	var cfg model.AutoTestConfig
	h.DB.Where("tenant_id = ?", tid).FirstOrCreate(&cfg, model.AutoTestConfig{TenantID: tid})
	frequencyChanged := cfg.Frequency != "" && cfg.Frequency != req.Frequency

	cfg.IsEnabled = req.IsEnabled
	cfg.Frequency = req.Frequency
	cfg.TargetMode = req.TargetMode
	cfg.SamplePercent = req.SamplePercent

	// Calculate next run if enabling
	if req.IsEnabled && (cfg.NextRunAt == nil || frequencyChanged) {
		next := nextRun(req.Frequency)
		cfg.NextRunAt = &next
	}
	if !req.IsEnabled {
		cfg.NextRunAt = nil
	}

	h.DB.Save(&cfg)
	c.JSON(http.StatusOK, cfg)
}

func nextRun(freq string) time.Time {
	now := time.Now().UTC()
	switch freq {
	case "monthly":
		return now.AddDate(0, 1, 0)
	case "quarterly":
		return now.AddDate(0, 3, 0)
	case "biannual":
		return now.AddDate(0, 6, 0)
	default:
		return now.AddDate(0, 3, 0)
	}
}
