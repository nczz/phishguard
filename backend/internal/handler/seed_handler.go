package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phishguard/phishguard/internal/middleware"
	"github.com/phishguard/phishguard/internal/service"
)

func (h *Handler) SeedSampleData(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	if err := service.SeedTenantData(h.DB, tid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "範例資料已匯入"})
}
