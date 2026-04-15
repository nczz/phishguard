package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phishguard/phishguard/internal/mailer"
	"github.com/phishguard/phishguard/internal/middleware"
	"github.com/phishguard/phishguard/internal/model"
)

func (h *Handler) ListSMTPProfiles(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	profiles, err := h.SMTPRepo.FindAllByTenant(tid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profiles)
}

func (h *Handler) CreateSMTPProfile(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	var req struct {
		Name          string `json:"name" binding:"required"`
		MailerType    string `json:"mailer_type"`
		Host          string `json:"host"`
		Port          *int   `json:"port"`
		Username      string `json:"username"`
		Password      string `json:"password"`
		FromAddress   string `json:"from_address"`
		FromName      string `json:"from_name"`
		TLSRequired   bool   `json:"tls_required"`
		MailgunDomain string `json:"mailgun_domain"`
		MailgunAPIKey string `json:"mailgun_api_key"`
		SESRegion     string `json:"ses_region"`
		SESAccessKey  string `json:"ses_access_key"`
		SESSecretKey  string `json:"ses_secret_key"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate from_address
	if !strings.Contains(req.FromAddress, "@") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "寄件地址必須是完整的 email 格式（如 noreply@yourdomain.com）"})
		return
	}
	if req.MailerType == "mailgun" && req.MailgunDomain != "" {
		fromDomain := req.FromAddress[strings.LastIndex(req.FromAddress, "@")+1:]
		if fromDomain != req.MailgunDomain {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("寄件地址的域名（%s）必須與 Mailgun Domain（%s）一致", fromDomain, req.MailgunDomain)})
			return
		}
	}
	p := model.SMTPProfile{
		TenantID:      tid,
		Name:          req.Name,
		MailerType:    req.MailerType,
		Host:          req.Host,
		Port:          req.Port,
		Username:      req.Username,
		PasswordEnc:   []byte(req.Password),
		FromAddress:   req.FromAddress,
		FromName:      req.FromName,
		TLSRequired:   req.TLSRequired,
		MailgunDomain: req.MailgunDomain,
		MailgunAPIKey: []byte(req.MailgunAPIKey),
		SESRegion:     req.SESRegion,
		SESAccessKey:  []byte(req.SESAccessKey),
		SESSecretKey:  []byte(req.SESSecretKey),
	}
	if err := h.SMTPRepo.Create(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *Handler) TestSMTPProfile(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	profile, err := h.SMTPRepo.FindByID(tid, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	var req struct {
		To string `json:"to" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := map[string]string{
		"host":       profile.Host,
		"port":       fmt.Sprintf("%d", derefInt(profile.Port)),
		"username":   profile.Username,
		"password":   string(profile.PasswordEnc),
		"tls":        fmt.Sprintf("%t", profile.TLSRequired),
		"domain":     profile.MailgunDomain,
		"api_key":    string(profile.MailgunAPIKey),
		"region":     profile.SESRegion,
		"access_key": string(profile.SESAccessKey),
		"secret_key": string(profile.SESSecretKey),
	}
	m, err := mailer.NewMailer(profile.MailerType, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	msg := &mailer.Message{
		From:     profile.FromAddress,
		FromName: profile.FromName,
		To:       req.To,
		Subject:  "PhishGuard SMTP Test",
		TextBody: "This is a test email.",
	}
	if err := m.Send(context.Background(), msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "test email sent"})
}

func derefInt(p *int) int {
	if p != nil {
		return *p
	}
	return 0
}
