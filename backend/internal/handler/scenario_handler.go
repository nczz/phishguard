package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/middleware"
	"github.com/nczz/phishguard/internal/model"
)

func (h *Handler) ListScenarios(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	scenarios, err := h.ScenarioRepo.FindAllByTenant(tid)
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, scenarios)
}

func (h *Handler) CreateScenario(c *gin.Context) {
	var req struct {
		Name          string `json:"name" binding:"required"`
		Category      string `json:"category" binding:"required"`
		Difficulty    int    `json:"difficulty"`
		Language      string `json:"language"`
		TemplateID    *int64 `json:"template_id"`
		PageID        *int64 `json:"page_id"`
		EducationHTML string `json:"education_html"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s := model.Scenario{
		TenantID: middleware.GetContextTenantID(c),
		Name: req.Name, Category: req.Category, Difficulty: req.Difficulty,
		Language: req.Language, TemplateID: req.TemplateID, PageID: req.PageID,
		EducationHTML: req.EducationHTML, IsActive: true,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if s.Language == "" { s.Language = "zh-TW" }
	if s.Difficulty == 0 { s.Difficulty = 2 }
	if err := h.ScenarioRepo.Create(&s); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, s)
}

func (h *Handler) UpdateScenario(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var s model.Scenario
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s.ID = id
	s.TenantID = &tid
	if err := h.ScenarioRepo.Update(tid, &s); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *Handler) DeleteScenario(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ScenarioRepo.Delete(tid, id); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
