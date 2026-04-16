package tracker

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/model"
	"gorm.io/gorm"
)

// 1x1 transparent GIF (GIF89a)
var transparentGIF = []byte{
	0x47, 0x49, 0x46, 0x38, 0x39, 0x61, // GIF89a
	0x01, 0x00, 0x01, 0x00, // 1x1
	0x80, 0x00, 0x00, // GCT flag, 1 color
	0xff, 0xff, 0xff, // color 0: white
	0x00, 0x00, 0x00, // color 1: black
	0x21, 0xf9, 0x04, // GCE
	0x01, 0x00, 0x00, 0x00, 0x00, // transparent index 0
	0x2c, 0x00, 0x00, 0x00, 0x00, // image descriptor
	0x01, 0x00, 0x01, 0x00, 0x00, // 1x1, no LCT
	0x02, 0x02, 0x4c, 0x01, 0x00, // LZW min code size 2, data
	0x3b, // trailer
}

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	t := r.Group("/t")
	t.GET("/o/:rid", h.HandleOpen)
	t.GET("/c/:rid", h.HandleClick)
	t.GET("/d/:rid/:filename", h.HandleDownload)
	t.POST("/s/:rid", h.HandleSubmit)
	t.POST("/r/:rid", h.HandleReport)
	t.GET("/r/:rid", h.HandleReport)
	t.GET("/landing", h.HandleLanding)
}

// serveHTML writes HTML with CSP headers that block all script execution.
// HTML content is NOT sanitized to preserve tenant's visual design.
// CSP script-src 'none' blocks inline scripts, event handlers, and javascript: URLs at browser level.
func serveHTML(c *gin.Context, html string) {
	c.Header("Content-Security-Policy", "script-src 'none'; frame-ancestors 'none'")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func (h *Handler) HandleOpen(c *gin.Context) {
	var result model.Result
	if err := h.DB.Where("rid = ?", c.Param("rid")).First(&result).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	now := time.Now()
	h.DB.Model(&result).Where("opened_at IS NULL").Update("opened_at", now)
	recordEvent(h.DB, result.ID, result.CampaignID, model.EventOpened, c.Request, nil)

	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "image/gif", transparentGIF)
}

func (h *Handler) HandleClick(c *gin.Context) {
	var result model.Result
	if err := h.DB.Where("rid = ?", c.Param("rid")).First(&result).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	now := time.Now()
	h.DB.Model(&result).Where("opened_at IS NULL").Update("opened_at", now)
	h.DB.Model(&result).Where("clicked_at IS NULL").Update("clicked_at", now)
	recordEvent(h.DB, result.ID, result.CampaignID, model.EventClicked, c.Request, nil)

	// Build landing URL from campaign's phish_url
	var campaign model.Campaign
	h.DB.First(&campaign, result.CampaignID)
	phishURL := strings.TrimRight(campaign.PhishURL, "/")
	c.Redirect(http.StatusFound, phishURL+"/t/landing?rid="+result.RID)
}

// HandleDownload does not exist as a constant — we don't need the Campaign preload here.
func (h *Handler) HandleDownload(c *gin.Context) {
	var result model.Result
	if err := h.DB.Where("rid = ?", c.Param("rid")).First(&result).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	now := time.Now()
	h.DB.Model(&result).Where("opened_at IS NULL").Update("opened_at", now)
	h.DB.Model(&result).Where("clicked_at IS NULL").Update("clicked_at", now)
	recordEvent(h.DB, result.ID, result.CampaignID, "downloaded", c.Request, map[string]interface{}{
		"filename": c.Param("filename"),
	})

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("<html><body><h1>File not available</h1></body></html>"))
}

func (h *Handler) HandleSubmit(c *gin.Context) {
	var result model.Result
	if err := h.DB.Where("rid = ?", c.Param("rid")).First(&result).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	now := time.Now()
	h.DB.Model(&result).Where("opened_at IS NULL").Update("opened_at", now)
	h.DB.Model(&result).Where("clicked_at IS NULL").Update("clicked_at", now)
	h.DB.Model(&result).Where("submitted_at IS NULL").Update("submitted_at", now)

	// Record field names only (not values)
	_ = c.Request.ParseForm()
	fields := make([]string, 0, len(c.Request.PostForm))
	for k := range c.Request.PostForm {
		fields = append(fields, k)
	}
	recordEvent(h.DB, result.ID, result.CampaignID, model.EventSubmitted, c.Request, map[string]interface{}{
		"fields": fields,
	})

	// Look up education HTML from scenario
	html := "<html><body><h1>Training Complete</h1><p>This was a phishing simulation.</p></body></html>"
	var campaign model.Campaign
	if h.DB.First(&campaign, result.CampaignID).Error == nil && campaign.ScenarioID != nil {
		var scenario model.Scenario
		if h.DB.First(&scenario, *campaign.ScenarioID).Error == nil && scenario.EducationHTML != "" {
			html = scenario.EducationHTML
		}
	}

	serveHTML(c, html)
}

func (h *Handler) HandleReport(c *gin.Context) {
	var result model.Result
	if err := h.DB.Where("rid = ?", c.Param("rid")).First(&result).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	now := time.Now()
	h.DB.Model(&result).Where("reported_at IS NULL").Update("reported_at", now)
	recordEvent(h.DB, result.ID, result.CampaignID, model.EventReported, c.Request, nil)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<!DOCTYPE html><html><head><meta charset="utf-8"><title>舉報成功</title>
<style>body{font-family:-apple-system,sans-serif;display:flex;justify-content:center;align-items:center;min-height:100vh;background:#f0f2f5;margin:0}
.card{background:#fff;padding:48px;border-radius:12px;box-shadow:0 2px 8px rgba(0,0,0,.1);text-align:center;max-width:480px}
h1{color:#52c41a;margin-bottom:16px}p{color:#666;line-height:1.8}</style></head>
<body><div class="card"><h1>✅ 感謝您的舉報！</h1>
<p>您已成功舉報這封可疑信件。<br>這是公司資安團隊發送的<strong>釣魚模擬測試</strong>。</p>
<p>您的警覺性非常好！能夠辨識並舉報可疑信件，是保護公司資安的重要行為。</p>
<p style="color:#999;font-size:13px;margin-top:24px;">本測試由 PhishGuard 釣魚模擬平台提供</p></div></body></html>`))
}

func (h *Handler) HandleLanding(c *gin.Context) {
	rid := c.Query("rid")
	if rid == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	var result model.Result
	if err := h.DB.Where("rid = ?", rid).First(&result).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	var campaign model.Campaign
	if err := h.DB.First(&campaign, result.CampaignID).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// Resolve landing page: campaign.PageID > scenario.PageID
	var pageID *int64
	if campaign.PageID != nil {
		pageID = campaign.PageID
	} else if campaign.ScenarioID != nil {
		var scenario model.Scenario
		if h.DB.First(&scenario, *campaign.ScenarioID).Error == nil {
			pageID = scenario.PageID
		}
	}

	if pageID == nil {
		c.Status(http.StatusNotFound)
		return
	}

	var page model.LandingPage
	if err := h.DB.First(&page, *pageID).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// Inject rid into form action URL
	phishURL := strings.TrimRight(campaign.PhishURL, "/")
	html := strings.ReplaceAll(page.HTML, "{{.RID}}", rid)
	html = strings.ReplaceAll(html, "{{.SubmitURL}}", phishURL+"/t/s/"+rid)

	serveHTML(c, html)
}

func recordEvent(db *gorm.DB, resultID, campaignID int64, eventType string, r *http.Request, detail map[string]interface{}) {
	var detailStr string
	if detail != nil {
		if b, err := json.Marshal(detail); err == nil {
			detailStr = string(b)
		}
	}
	db.Create(&model.Event{
		ResultID:   resultID,
		CampaignID: campaignID,
		EventType:  eventType,
		IPAddress:  r.RemoteAddr,
		UserAgent:  r.UserAgent(),
		Detail:     detailStr,
	})
}
