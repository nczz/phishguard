package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/middleware"
	"github.com/nczz/phishguard/internal/model"
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

	var cfg model.AutoTestConfig
	h.DB.Where("tenant_id = ?", tid).FirstOrCreate(&cfg, model.AutoTestConfig{TenantID: tid})

	cfg.IsEnabled = req.IsEnabled
	cfg.Frequency = req.Frequency
	cfg.TargetMode = req.TargetMode
	cfg.SamplePercent = req.SamplePercent

	// Calculate next run if enabling
	if req.IsEnabled && cfg.NextRunAt == nil {
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
	now := time.Now()
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
