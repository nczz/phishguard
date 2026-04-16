package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/middleware"
	"github.com/nczz/phishguard/internal/model"
)

// PauseCampaign — stops sending but keeps scheduled results for later resume
func (h *Handler) PauseCampaign(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var campaign model.Campaign
	if err := h.DB.Where("id = ? AND tenant_id = ?", id, tid).First(&campaign).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
		return
	}
	if campaign.Status != "sending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只有發送中的活動可以暫停"})
		return
	}

	h.DB.Model(&campaign).Update("status", "paused")
	c.JSON(http.StatusOK, gin.H{"message": "活動已暫停，未發送的信件將保留排程"})
}

// ResumeCampaign — resume a paused campaign
func (h *Handler) ResumeCampaign(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var campaign model.Campaign
	if err := h.DB.Where("id = ? AND tenant_id = ?", id, tid).First(&campaign).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
		return
	}
	if campaign.Status != "paused" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只有已暫停的活動可以恢復"})
		return
	}

	h.DB.Model(&campaign).Update("status", "sending")
	c.JSON(http.StatusOK, gin.H{"message": "活動已恢復發送"})
}

// StopCampaign — permanently stop, mark remaining as cancelled
func (h *Handler) StopCampaign(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var campaign model.Campaign
	if err := h.DB.Where("id = ? AND tenant_id = ?", id, tid).First(&campaign).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
		return
	}
	if campaign.Status != "sending" && campaign.Status != "paused" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只有發送中或已暫停的活動可以終止"})
		return
	}

	// Mark remaining scheduled results as cancelled
	h.DB.Model(&model.Result{}).
		Where("campaign_id = ? AND status = ?", id, "scheduled").
		Updates(map[string]interface{}{"status": "cancelled", "error_detail": "活動已被手動終止"})

	now := time.Now()
	h.DB.Model(&campaign).Updates(map[string]interface{}{"status": "stopped", "completed_at": now})

	// Count what happened
	var sent, cancelled int64
	h.DB.Model(&model.Result{}).Where("campaign_id = ? AND status = ?", id, "sent").Count(&sent)
	h.DB.Model(&model.Result{}).Where("campaign_id = ? AND status = ?", id, "cancelled").Count(&cancelled)

	c.JSON(http.StatusOK, gin.H{
		"message":   "活動已終止",
		"sent":      sent,
		"cancelled": cancelled,
	})
}
