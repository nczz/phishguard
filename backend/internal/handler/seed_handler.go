package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/middleware"
	"github.com/nczz/phishguard/internal/service"
)

func (h *Handler) SeedSampleData(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	if err := service.SeedTenantData(h.DB, tid); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "範例資料已匯入"})
}
