package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/middleware"
	"github.com/nczz/phishguard/internal/model"
	"gorm.io/gorm"
)

// --- Tenant Dashboard with real stats ---

type DashboardStats struct {
	TotalCampaigns int64   `json:"total_campaigns"`
	AvgOpenRate    float64 `json:"avg_open_rate"`
	AvgClickRate   float64 `json:"avg_click_rate"`
	AvgSubmitRate  float64 `json:"avg_submit_rate"`
	AvgReportRate  float64 `json:"avg_report_rate"`
}

func (h *Handler) TenantDashboardStats(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	var campaigns []model.Campaign
	h.DB.Where("tenant_id = ? AND status IN ?", tid, []string{"sending", "sent", "completed"}).Find(&campaigns)

	stats := DashboardStats{TotalCampaigns: int64(len(campaigns))}
	if len(campaigns) == 0 {
		c.JSON(http.StatusOK, stats)
		return
	}

	var totalOpen, totalClick, totalSubmit, totalReport float64
	for _, camp := range campaigns {
		f, err := h.ResultRepo.GetFunnelStats(tid, camp.ID)
		if err != nil || f.Total == 0 {
			continue
		}
		totalOpen += float64(f.Opened) / float64(f.Total)
		totalClick += float64(f.Clicked) / float64(f.Total)
		totalSubmit += float64(f.Submitted) / float64(f.Total)
		totalReport += float64(f.Reported) / float64(f.Total)
	}
	n := float64(len(campaigns))
	stats.AvgOpenRate = round2(totalOpen / n * 100)
	stats.AvgClickRate = round2(totalClick / n * 100)
	stats.AvgSubmitRate = round2(totalSubmit / n * 100)
	stats.AvgReportRate = round2(totalReport / n * 100)
	c.JSON(http.StatusOK, stats)
}

// --- Campaign recipient detail list ---

type RecipientResult struct {
	Email       string     `json:"email"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Department  string     `json:"department"`
	Status      string     `json:"status"`
	SentAt      *time.Time `json:"sent_at"`
	OpenedAt    *time.Time `json:"opened_at"`
	ClickedAt   *time.Time `json:"clicked_at"`
	SubmittedAt *time.Time `json:"submitted_at"`
	ReportedAt  *time.Time `json:"reported_at"`
}

func (h *Handler) CampaignRecipients(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	cid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var results []model.Result
	h.DB.Where("campaign_id = ? AND tenant_id = ?", cid, tid).Preload("Recipient").Find(&results)

	list := make([]RecipientResult, 0, len(results))
	for _, r := range results {
		rr := RecipientResult{Status: r.Status, SentAt: r.SentAt, OpenedAt: r.OpenedAt, ClickedAt: r.ClickedAt, SubmittedAt: r.SubmittedAt, ReportedAt: r.ReportedAt}
		if r.Recipient != nil {
			rr.Email = r.Recipient.Email
			rr.FirstName = r.Recipient.FirstName
			rr.LastName = r.Recipient.LastName
			rr.Department = r.Recipient.Department
		}
		list = append(list, rr)
	}
	c.JSON(http.StatusOK, list)
}

// --- Campaign CSV export ---

func (h *Handler) ExportCampaignCSV(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	cid, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var campaign model.Campaign
	if err := h.DB.Where("id = ? AND tenant_id = ?", cid, tid).First(&campaign).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	var results []model.Result
	h.DB.Where("campaign_id = ? AND tenant_id = ?", cid, tid).Preload("Recipient").Find(&results)

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s_results.csv", campaign.Name))
	// UTF-8 BOM for Excel
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	w := csv.NewWriter(c.Writer)
	w.Write([]string{"email", "last_name", "first_name", "department", "status", "sent_at", "opened_at", "clicked_at", "submitted_at", "reported_at"})
	for _, r := range results {
		row := []string{"", "", "", "", r.Status, fmtTime(r.SentAt), fmtTime(r.OpenedAt), fmtTime(r.ClickedAt), fmtTime(r.SubmittedAt), fmtTime(r.ReportedAt)}
		if r.Recipient != nil {
			row[0] = r.Recipient.Email
			row[1] = r.Recipient.LastName
			row[2] = r.Recipient.FirstName
			row[3] = r.Recipient.Department
		}
		w.Write(row)
	}
	w.Flush()
}

// --- Repeat offender tracking ---

type OffenderRecord struct {
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Department string `json:"department"`
	History    []struct {
		CampaignID   int64  `json:"campaign_id"`
		CampaignName string `json:"campaign_name"`
		Status       string `json:"status"`
	} `json:"history"`
	ClickCount  int `json:"click_count"`
	SubmitCount int `json:"submit_count"`
}

func (h *Handler) RepeatOffenders(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)

	// Get all results with recipient and campaign info
	var results []model.Result
	h.DB.Where("tenant_id = ?", tid).Preload("Recipient").Find(&results)

	// Load campaign names
	var campaigns []model.Campaign
	h.DB.Where("tenant_id = ?", tid).Find(&campaigns)
	campNames := map[int64]string{}
	for _, c := range campaigns {
		campNames[c.ID] = c.Name
	}

	// Group by email
	type entry struct {
		r    *model.Recipient
		hist []struct {
			CampaignID   int64  `json:"campaign_id"`
			CampaignName string `json:"campaign_name"`
			Status       string `json:"status"`
		}
		clicks, submits int
	}
	byEmail := map[string]*entry{}
	for _, r := range results {
		if r.Recipient == nil {
			continue
		}
		email := r.Recipient.Email
		e, ok := byEmail[email]
		if !ok {
			e = &entry{r: r.Recipient}
			byEmail[email] = e
		}
		status := bestStatus(r)
		e.hist = append(e.hist, struct {
			CampaignID   int64  `json:"campaign_id"`
			CampaignName string `json:"campaign_name"`
			Status       string `json:"status"`
		}{r.CampaignID, campNames[r.CampaignID], status})
		if r.ClickedAt != nil {
			e.clicks++
		}
		if r.SubmittedAt != nil {
			e.submits++
		}
	}

	// Filter: only people who clicked or submitted at least once
	offenders := []OffenderRecord{}
	for _, e := range byEmail {
		if e.clicks == 0 && e.submits == 0 {
			continue
		}
		offenders = append(offenders, OffenderRecord{
			Email: e.r.Email, FirstName: e.r.FirstName, LastName: e.r.LastName, Department: e.r.Department,
			History: e.hist, ClickCount: e.clicks, SubmitCount: e.submits,
		})
	}
	c.JSON(http.StatusOK, offenders)
}

// --- Trend analysis ---

type TrendPoint struct {
	CampaignID   int64   `json:"campaign_id"`
	CampaignName string  `json:"campaign_name"`
	LaunchedAt   string  `json:"launched_at"`
	OpenRate     float64 `json:"open_rate"`
	ClickRate    float64 `json:"click_rate"`
	SubmitRate   float64 `json:"submit_rate"`
	ReportRate   float64 `json:"report_rate"`
}

func (h *Handler) TrendAnalysis(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	var campaigns []model.Campaign
	h.DB.Where("tenant_id = ? AND status IN ?", tid, []string{"sending", "sent", "completed"}).Order("launched_at ASC").Find(&campaigns)

	points := make([]TrendPoint, 0, len(campaigns))
	for _, camp := range campaigns {
		f, err := h.ResultRepo.GetFunnelStats(tid, camp.ID)
		if err != nil || f.Total == 0 {
			continue
		}
		launched := ""
		if camp.LaunchedAt != nil {
			launched = camp.LaunchedAt.Format("2006-01-02")
		}
		points = append(points, TrendPoint{
			CampaignID: camp.ID, CampaignName: camp.Name, LaunchedAt: launched,
			OpenRate:   round2(float64(f.Opened) / float64(f.Total) * 100),
			ClickRate:  round2(float64(f.Clicked) / float64(f.Total) * 100),
			SubmitRate: round2(float64(f.Submitted) / float64(f.Total) * 100),
			ReportRate: round2(float64(f.Reported) / float64(f.Total) * 100),
		})
	}
	c.JSON(http.StatusOK, points)
}

// --- Helpers ---

func round2(f float64) float64 { return float64(int(f*100)) / 100 }

func fmtTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func bestStatus(r model.Result) string {
	if r.SubmittedAt != nil {
		return "submitted"
	}
	if r.ClickedAt != nil {
		return "clicked"
	}
	if r.OpenedAt != nil {
		return "opened"
	}
	if r.SentAt != nil {
		return "sent"
	}
	return r.Status
}

// Ensure gorm is used
var _ *gorm.DB
