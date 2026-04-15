package handler

import (
	"net/http"
	"strconv"

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
			Gender     string `json:"gender"`
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
	recipients := make([]model.Recipient, 0, len(req.Recipients))
	seen := make(map[string]bool)
	for _, r := range req.Recipients {
		email := r.Email
		if email == "" || seen[email] {
			continue
		}
		seen[email] = true
		recipients = append(recipients, model.Recipient{
			TenantID:   tid,
			GroupID:    req.GroupID,
			Email:      email,
			FirstName:  r.FirstName,
			LastName:   r.LastName,
			Department: r.Department,
			Gender:     r.Gender,
			Position:   r.Position,
		})
	}
	created, updated, err := h.RecipientRepo.UpsertRecipients(tid, req.GroupID, recipients)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"created": created, "updated": updated, "total": created + updated})
}

func (h *Handler) UpdateRecipient(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Email      string `json:"email"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Department string `json:"department"`
		Gender     string `json:"gender"`
		Position   string `json:"position"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.RecipientRepo.UpdateRecipient(tid, id, req.Email, req.FirstName, req.LastName, req.Department, req.Gender, req.Position); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) DeleteRecipient(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.RecipientRepo.DeleteRecipient(tid, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
