package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/mailer"
	"github.com/nczz/phishguard/internal/middleware"
	"github.com/nczz/phishguard/internal/service"
)

func (h *Handler) SendCampaignReportEmail(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid campaign id"})
		return
	}

	campaign, err := h.CampaignRepo.FindByID(tid, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
		return
	}

	profile, err := h.SMTPRepo.FindByID(tid, campaign.SMTPProfileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "smtp profile not found"})
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
		serverError(c, err)
		return
	}

	if err := service.SendCampaignReport(h.DB, h.ResultRepo, campaign, m, profile.FromAddress); err != nil {
		serverError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "report sent"})
}
