package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phishguard/phishguard/internal/middleware"
	"github.com/phishguard/phishguard/internal/model"
	"github.com/phishguard/phishguard/internal/service"
)

func (h *Handler) CreateCampaign(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)

	// Plan limit check: campaigns per year
	var tenant model.Tenant
	if err := h.DB.First(&tenant, tid).Error; err == nil {
		limits := service.GetEffectiveLimits(&tenant)
		if limits.MaxCampaignsPerYear > 0 {
			count, _ := h.CampaignRepo.CountByTenantThisYear(tid)
			if count >= int64(limits.MaxCampaignsPerYear) {
				c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("已達年度活動上限 (%d 次)，請升級方案", limits.MaxCampaignsPerYear)})
				return
			}
		}
	}

	var req service.CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	campaign, err := h.CampaignService.CreateCampaign(tid, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, campaign)
}

func (h *Handler) ListCampaigns(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	campaigns, err := h.CampaignRepo.FindAllByTenant(tid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, campaigns)
}

func (h *Handler) GetCampaign(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	campaign, err := h.CampaignRepo.FindByID(tid, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
		return
	}
	c.JSON(http.StatusOK, campaign)
}

func (h *Handler) LaunchCampaign(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Plan limit check: monthly email quota
	var tenant model.Tenant
	if err := h.DB.First(&tenant, tid).Error; err == nil {
		limits := service.GetEffectiveLimits(&tenant)
		if limits.MaxEmailsPerMonth > 0 {
			sent, _ := h.ResultRepo.CountSentThisMonth(tid)
			// Estimate: count recipients in this campaign
			var recipientCount int64
			h.DB.Model(&model.Result{}).Where("campaign_id = ?", id).Count(&recipientCount)
			if sent+recipientCount > int64(limits.MaxEmailsPerMonth) {
				c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("本月發信量將超過上限 (%d 封)，已發送 %d 封，本次需發送 %d 封。請升級方案或聯繫管理員。", limits.MaxEmailsPerMonth, sent, recipientCount)})
				return
			}
		}
	}

	if err := h.CampaignService.LaunchCampaign(tid, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "campaign launched"})
}

func (h *Handler) DeleteCampaign(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.CampaignRepo.Delete(tid, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
