package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/middleware"
	"github.com/nczz/phishguard/internal/model"
	"github.com/nczz/phishguard/internal/service"
)

func (h *Handler) ListRecipientGroups(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	groups, err := h.RecipientRepo.FindGroupsByTenant(tid)
	if err != nil {
		serverError(c, err)
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
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, g)
}

func (h *Handler) ImportRecipients(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)

	var req struct {
		GroupID    int64 `json:"group_id" binding:"required"`
		Sync       bool  `json:"sync"`
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

	var tenant model.Tenant
	if err := h.DB.First(&tenant, tid).Error; err == nil {
		limits := service.GetEffectiveLimits(&tenant)
		if limits.MaxRecipients > 0 {
			current, _ := h.RecipientRepo.CountByTenant(tid)
			emails := make([]string, 0, len(recipients))
			for _, r := range recipients {
				emails = append(emails, r.Email)
			}
			var existing int64
			if len(emails) > 0 {
				h.DB.Model(&model.Recipient{}).Where("tenant_id = ? AND email IN ?", tid, emails).Count(&existing)
			}
			toCreate := int64(len(recipients)) - existing
			if current+toCreate > int64(limits.MaxRecipients) {
				c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("匯入後將超過收件人上限 (%d 人)，目前 %d 人，本次新增 %d 人", limits.MaxRecipients, current, toCreate)})
				return
			}
		}
	}

	created, updated, deactivated, err := h.RecipientRepo.UpsertRecipients(tid, req.GroupID, recipients, req.Sync)
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"created": created, "updated": updated, "deactivated": deactivated, "total": created + updated})
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
		serverError(c, err)
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
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) BatchDeleteRecipients(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	var req struct {
		IDs []int64 `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.RecipientRepo.BatchDelete(tid, req.IDs); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": len(req.IDs)})
}

func (h *Handler) BatchSetActiveRecipients(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	var req struct {
		IDs    []int64 `json:"ids" binding:"required"`
		Active bool    `json:"active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.RecipientRepo.BatchSetActive(tid, req.IDs, req.Active); err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": len(req.IDs), "active": req.Active})
}

func (h *Handler) ValidateRecipients(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)

	// Get all recipients for this tenant
	count, _ := h.RecipientRepo.CountByTenant(tid)
	var recipients []model.Recipient
	h.DB.Where("tenant_id = ?", tid).Find(&recipients)

	type FilteredItem struct {
		Email  string `json:"email"`
		Name   string `json:"name"`
		Reason string `json:"reason"`
	}

	filtered := []FilteredItem{}
	valid := 0
	for _, r := range recipients {
		reason := service.GetEmailFilterReason(r.Email)
		if reason != "" {
			filtered = append(filtered, FilteredItem{
				Email:  r.Email,
				Name:   r.LastName + r.FirstName,
				Reason: reason,
			})
		} else {
			valid++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    count,
		"valid":    valid,
		"filtered": filtered,
	})
}
