package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phishguard/phishguard/internal/middleware"
	"github.com/phishguard/phishguard/internal/model"
	"github.com/phishguard/phishguard/internal/service"
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

	// Plan limit check: max recipients
	var tenant model.Tenant
	if err := h.DB.First(&tenant, tid).Error; err == nil {
		limits := service.GetEffectiveLimits(&tenant)
		if limits.MaxRecipients > 0 {
			current, _ := h.RecipientRepo.CountByTenant(tid)
			if current >= int64(limits.MaxRecipients) {
				c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("已達收件人上限 (%d 人)，請升級方案或聯繫管理員調整", limits.MaxRecipients)})
				return
			}
		}
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": len(req.IDs)})
}

func (h *Handler) ValidateRecipients(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	var results []model.Result
	h.DB.Where("tenant_id = ?", tid).Find(&results) // just to use tid

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
		reason := getFilterReason(r.Email)
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

func getFilterReason(email string) string {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return "無效的 email 格式"
	}
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "無效的 email 格式"
	}
	domain := strings.ToLower(parts[1])

	blockedDomains := map[string]string{
		"gmail.com": "公共信箱（Gmail）", "yahoo.com": "公共信箱（Yahoo）",
		"hotmail.com": "公共信箱（Hotmail）", "outlook.com": "公共信箱（Outlook）",
		"aol.com": "公共信箱（AOL）", "icloud.com": "公共信箱（iCloud）",
		"mail.com": "公共信箱", "protonmail.com": "公共信箱（ProtonMail）",
		"zoho.com": "公共信箱（Zoho）", "yandex.com": "公共信箱（Yandex）",
		"gmx.com": "公共信箱（GMX）", "live.com": "公共信箱（Live）",
	}
	if reason, ok := blockedDomains[domain]; ok {
		return reason + " — 釣魚測試僅限企業域名"
	}

	roleAddresses := map[string]string{
		"abuse@": "角色信箱（abuse）", "postmaster@": "角色信箱（postmaster）",
		"hostmaster@": "角色信箱", "webmaster@": "角色信箱",
		"noc@": "角色信箱", "security@": "角色信箱",
		"mailer-daemon@": "系統信箱", "admin@": "角色信箱（admin）",
	}
	lower := strings.ToLower(email)
	for prefix, reason := range roleAddresses {
		if strings.HasPrefix(lower, prefix) {
			return reason + " — 不應作為釣魚測試對象"
		}
	}

	return ""
}
