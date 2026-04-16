package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/middleware"
	"github.com/nczz/phishguard/internal/model"
)

func (h *Handler) ListPages(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	pages, err := h.PageRepo.FindAllByTenant(tid)
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, pages)
}

func (h *Handler) CreatePage(c *gin.Context) {
	var req struct {
		Name               string `json:"name" binding:"required"`
		HTML               string `json:"html" binding:"required"`
		CaptureCredentials bool   `json:"capture_credentials"`
		CaptureFields      string `json:"capture_fields"`
		RedirectURL        string `json:"redirect_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := model.LandingPage{
		TenantID: middleware.GetContextTenantID(c),
		Name: req.Name, HTML: req.HTML,
		CaptureCredentials: req.CaptureCredentials,
		CaptureFields: req.CaptureFields, RedirectURL: req.RedirectURL,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := h.PageRepo.Create(&p); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *Handler) UpdatePage(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var p model.LandingPage
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = id
	if err := h.PageRepo.Update(tid, &p); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *Handler) DeletePage(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.PageRepo.Delete(tid, id); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
