package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/middleware"
)

func (h *Handler) ListAuditLogs(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	logs, total, err := h.AuditRepo.FindByTenant(tid, limit, offset)
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs, "total": total})
}
