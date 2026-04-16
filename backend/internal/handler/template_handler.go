package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/middleware"
	"github.com/nczz/phishguard/internal/model"
)

func (h *Handler) ListTemplates(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	templates, err := h.TemplateRepo.FindAllByTenant(tid)
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, templates)
}

func (h *Handler) CreateTemplate(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Subject  string `json:"subject" binding:"required"`
		HTMLBody string `json:"html_body"`
		TextBody string `json:"text_body"`
		Category string `json:"category"`
		Language string `json:"language"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := middleware.GetUserID(c)
	t := model.EmailTemplate{
		TenantID:  middleware.GetContextTenantID(c),
		Name:      req.Name,
		Subject:   req.Subject,
		HTMLBody:  req.HTMLBody,
		TextBody:  req.TextBody,
		Category:  req.Category,
		Language:  req.Language,
		CreatedBy: &uid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if t.Language == "" {
		t.Language = "zh-TW"
	}
	if err := h.TemplateRepo.Create(&t); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *Handler) UpdateTemplate(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var t model.EmailTemplate
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t.ID = id
	if err := h.TemplateRepo.Update(tid, &t); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *Handler) DeleteTemplate(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.TemplateRepo.Delete(tid, id); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
