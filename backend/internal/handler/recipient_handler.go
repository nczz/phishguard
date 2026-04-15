package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phishguard/phishguard/internal/middleware"
	"github.com/phishguard/phishguard/internal/model"
)

func (h *Handler) ListRecipientGroups(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	groups, err := h.RecipientRepo.FindGroupsByTenant(tid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, groups)
}

func (h *Handler) CreateRecipientGroup(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g := model.RecipientGroup{TenantID: tid, Name: req.Name}
	if err := h.RecipientRepo.CreateGroup(&g); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, g)
}

func (h *Handler) ImportRecipients(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	var req struct {
		GroupID    int64 `json:"group_id" binding:"required"`
		Recipients []struct {
			Email      string `json:"email"`
			FirstName  string `json:"first_name"`
			LastName   string `json:"last_name"`
			Department string `json:"department"`
			Position   string `json:"position"`
		} `json:"recipients" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := h.RecipientRepo.FindGroupByID(tid, req.GroupID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}
	recipients := make([]model.Recipient, len(req.Recipients))
	for i, r := range req.Recipients {
		recipients[i] = model.Recipient{
			TenantID:   tid,
			GroupID:    req.GroupID,
			Email:      r.Email,
			FirstName:  r.FirstName,
			LastName:   r.LastName,
			Department: r.Department,
			Position:   r.Position,
		}
	}
	if err := h.RecipientRepo.BulkCreateRecipients(recipients); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"imported": len(recipients)})
}
