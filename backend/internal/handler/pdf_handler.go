package handler

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nczz/phishguard/internal/middleware"
	"github.com/nczz/phishguard/internal/model"
)

func (h *Handler) ExportCampaignPDFReal(c *gin.Context) {
	tid := *middleware.GetContextTenantID(c)
	cid, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var campaign model.Campaign
	if err := h.DB.Where("id = ? AND tenant_id = ?", cid, tid).First(&campaign).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	funnel, _ := h.ResultRepo.GetFunnelStats(tid, cid)
	depts, _ := h.ResultRepo.GetDepartmentStats(tid, cid)

	sort.Slice(depts, func(i, j int) bool {
		ri := float64(0)
		rj := float64(0)
		if depts[i].Total > 0 { ri = float64(depts[i].Clicked) / float64(depts[i].Total) }
		if depts[j].Total > 0 { rj = float64(depts[j].Clicked) / float64(depts[j].Total) }
		return ri > rj
	})

	pct := func(n int64) string {
		if funnel.Total == 0 { return "0%" }
		return fmt.Sprintf("%.1f%%", float64(n)/float64(funnel.Total)*100)
	}

	launched := "—"
	if campaign.LaunchedAt != nil {
		launched = campaign.LaunchedAt.Format("2006-01-02 15:04")
	}

	// Build department rows
	var deptRows strings.Builder
	for _, d := range depts {
		rate := 0.0
		if d.Total > 0 { rate = float64(d.Clicked) / float64(d.Total) * 100 }
		deptRows.WriteString(fmt.Sprintf(`<tr><td>%s</td><td>%d</td><td>%d</td><td>%.1f%%</td></tr>`, d.Department, d.Total, d.Clicked, rate))
	}

	html := fmt.Sprintf(`<!DOCTYPE html><html><head><meta charset="utf-8"><title>%s — 釣魚測試報告</title>
<style>
@media print { body { -webkit-print-color-adjust: exact; print-color-adjust: exact; } }
body { font-family: -apple-system, "Microsoft JhengHei", "PingFang TC", sans-serif; max-width: 800px; margin: 0 auto; padding: 40px; color: #333; }
h1 { text-align: center; color: #1677ff; margin-bottom: 4px; }
.subtitle { text-align: center; color: #999; margin-bottom: 32px; }
.info { background: #f5f5f5; padding: 16px; border-radius: 8px; margin-bottom: 24px; }
.info span { margin-right: 24px; }
table { width: 100%%; border-collapse: collapse; margin-bottom: 24px; }
th, td { border: 1px solid #e8e8e8; padding: 10px 14px; text-align: left; }
th { background: #fafafa; font-weight: 600; }
.funnel td:first-child { font-weight: 500; }
.bar { height: 20px; border-radius: 4px; display: inline-block; }
h2 { color: #1677ff; border-bottom: 2px solid #1677ff; padding-bottom: 6px; margin-top: 32px; }
.footer { text-align: center; color: #bbb; margin-top: 40px; font-size: 12px; }
</style></head><body>
<h1>🛡️ PhishGuard 釣魚測試報告</h1>
<p class="subtitle">%s</p>

<div class="info">
  <span><strong>活動名稱：</strong>%s</span>
  <span><strong>狀態：</strong>%s</span>
  <span><strong>發送時間：</strong>%s</span>
  <span><strong>收件人數：</strong>%d</span>
</div>

<h2>📊 釣魚漏斗</h2>
<table class="funnel">
<tr><th>指標</th><th>人數</th><th>比率</th><th>圖示</th></tr>
<tr><td>寄達</td><td>%d</td><td>%s</td><td><div class="bar" style="width:%dpx;background:#1677ff">&nbsp;</div></td></tr>
<tr><td>開信</td><td>%d</td><td>%s</td><td><div class="bar" style="width:%dpx;background:#13c2c2">&nbsp;</div></td></tr>
<tr><td>點擊</td><td>%d</td><td>%s</td><td><div class="bar" style="width:%dpx;background:#faad14">&nbsp;</div></td></tr>
<tr><td>下載</td><td>%d</td><td>%s</td><td><div class="bar" style="width:%dpx;background:#722ed1">&nbsp;</div></td></tr>
<tr><td>提交</td><td>%d</td><td>%s</td><td><div class="bar" style="width:%dpx;background:#ff4d4f">&nbsp;</div></td></tr>
<tr><td>舉報</td><td>%d</td><td>%s</td><td><div class="bar" style="width:%dpx;background:#52c41a">&nbsp;</div></td></tr>
</table>

<h2>🏢 部門風險排名</h2>
<table>
<tr><th>部門</th><th>總人數</th><th>點擊人數</th><th>點擊率</th></tr>
%s
</table>

<p class="footer">由 PhishGuard 釣魚模擬測試平台產生 · %s</p>
</body></html>`,
		campaign.Name,
		time.Now().Format("2006-01-02"),
		campaign.Name, campaign.Status, launched, funnel.Total,
		funnel.Sent, pct(funnel.Sent), barWidth(funnel.Sent, funnel.Total),
		funnel.Opened, pct(funnel.Opened), barWidth(funnel.Opened, funnel.Total),
		funnel.Clicked, pct(funnel.Clicked), barWidth(funnel.Clicked, funnel.Total),
		funnel.Downloaded, pct(funnel.Downloaded), barWidth(funnel.Downloaded, funnel.Total),
		funnel.Submitted, pct(funnel.Submitted), barWidth(funnel.Submitted, funnel.Total),
		funnel.Reported, pct(funnel.Reported), barWidth(funnel.Reported, funnel.Total),
		deptRows.String(),
		time.Now().Format("2006-01-02 15:04:05"),
	)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func barWidth(n, total int64) int {
	if total == 0 { return 0 }
	w := int(float64(n) / float64(total) * 300)
	if w < 2 && n > 0 { w = 2 }
	return w
}
