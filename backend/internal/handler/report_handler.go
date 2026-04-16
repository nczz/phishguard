package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/middleware"
)

func (h *Handler) GetCampaignReport(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid campaign id"})
		return
	}
	funnel, err := h.ReportService.GetCampaignFunnel(tid, id)
	if err != nil {
		serverError(c, err)
		return
	}
	depts, err := h.ReportService.GetDepartmentStats(tid, id)
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"funnel": funnel, "departments": depts})
}

func (h *Handler) ExportCampaignPDF(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "PDF export coming in Phase 2"})
}

func (h *Handler) GetOverviewReport(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	campaigns, err := h.CampaignRepo.FindAllByTenant(tid)
	if err != nil {
		serverError(c, err)
		return
	}
	ids := make([]int64, len(campaigns))
	for i, camp := range campaigns {
		ids[i] = camp.ID
	}
	overview, err := h.ReportService.GetOverview(tid, ids)
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, overview)
}

func (h *Handler) GetDepartmentReport(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	cidStr := c.Query("campaign_id")
	if cidStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "campaign_id is required"})
		return
	}
	cid, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid campaign_id"})
		return
	}
	stats, err := h.ReportService.GetDepartmentStats(tid, cid)
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, stats)
}
