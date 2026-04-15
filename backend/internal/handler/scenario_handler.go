package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phishguard/phishguard/internal/middleware"
	"github.com/phishguard/phishguard/internal/model"
)

func (h *Handler) ListScenarios(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	scenarios, err := h.ScenarioRepo.FindAllByTenant(tid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, scenarios)
}

func (h *Handler) CreateScenario(c *gin.Context) {
	var s model.Scenario
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.TenantID = middleware.GetContextTenantID(c)
	if err := h.ScenarioRepo.Create(&s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, s)
}
