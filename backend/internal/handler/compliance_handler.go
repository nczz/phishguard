package handler

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phishguard/phishguard/internal/middleware"
)

type ComplianceCheck struct {
	Name    string `json:"name"`
	Status  string `json:"status"` // pass / warn / fail
	Detail  string `json:"detail"`
	Fix     string `json:"fix,omitempty"`
}

type ComplianceResult struct {
	Domain string            `json:"domain"`
	Score  int               `json:"score"` // 0-100
	Checks []ComplianceCheck `json:"checks"`
}

func (h *Handler) CheckMailCompliance(c *gin.Context) {
	_ = middleware.GetContextTenantID(c)
	var req struct {
		FromAddress string `json:"from_address" binding:"required"`
		SmtpHost    string `json:"smtp_host"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parts := strings.SplitN(req.FromAddress, "@", 2)
	if len(parts) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}
	domain := parts[1]

	checks := []ComplianceCheck{}
	score := 100

	// 1. SPF
	checks = append(checks, checkSPF(domain, &score))

	// 2. DKIM (check for common selectors)
	checks = append(checks, checkDKIM(domain, &score))

	// 3. DMARC
	checks = append(checks, checkDMARC(domain, &score))

	// 4. MX record
	checks = append(checks, checkMX(domain, &score))

	// 5. Reverse DNS (if SMTP host provided)
	if req.SmtpHost != "" {
		checks = append(checks, checkReverseDNS(req.SmtpHost, &score))
	}

	// 6. Rate limiting advice
	checks = append(checks, ComplianceCheck{
		Name:   "發信速率控制",
		Status: "pass",
		Detail: "系統已內建速率限制（每秒最多 10 封），避免觸發收件伺服器的速率保護",
	})

	// 7. List-Unsubscribe header
	checks = append(checks, ComplianceCheck{
		Name:   "List-Unsubscribe 標頭",
		Status: "pass",
		Detail: "系統自動在每封信加入 List-Unsubscribe 標頭，符合 RFC 8058 規範",
	})

	// 8. Message-ID
	checks = append(checks, ComplianceCheck{
		Name:   "Message-ID 標頭",
		Status: "pass",
		Detail: "系統自動產生唯一 Message-ID，格式為 <rid@domain>",
	})

	if score < 0 {
		score = 0
	}

	c.JSON(http.StatusOK, ComplianceResult{Domain: domain, Score: score, Checks: checks})
}

func checkSPF(domain string, score *int) ComplianceCheck {
	txts, err := net.LookupTXT(domain)
	if err != nil {
		*score -= 25
		return ComplianceCheck{Name: "SPF 記錄", Status: "fail", Detail: "無法查詢 DNS TXT 記錄: " + err.Error(),
			Fix: fmt.Sprintf("在 %s 的 DNS 加入 TXT 記錄：v=spf1 include:_spf.google.com ~all（依實際發信來源調整）", domain)}
	}
	for _, txt := range txts {
		if strings.HasPrefix(txt, "v=spf1") {
			return ComplianceCheck{Name: "SPF 記錄", Status: "pass", Detail: "已設定: " + txt}
		}
	}
	*score -= 25
	return ComplianceCheck{Name: "SPF 記錄", Status: "fail", Detail: "未找到 SPF 記錄",
		Fix: fmt.Sprintf("在 %s 的 DNS 加入 TXT 記錄：v=spf1 include:mailgun.org include:amazonses.com ~all", domain)}
}

func checkDKIM(domain string, score *int) ComplianceCheck {
	selectors := []string{"default", "google", "selector1", "selector2", "k1", "mandrill", "mailgun", "ses"}
	for _, sel := range selectors {
		host := sel + "._domainkey." + domain
		txts, err := net.LookupTXT(host)
		if err == nil && len(txts) > 0 {
			for _, txt := range txts {
				if strings.Contains(txt, "p=") {
					return ComplianceCheck{Name: "DKIM 記錄", Status: "pass", Detail: fmt.Sprintf("已設定 (selector: %s)", sel)}
				}
			}
		}
	}
	*score -= 20
	return ComplianceCheck{Name: "DKIM 記錄", Status: "warn", Detail: "未找到常見 DKIM selector（已檢查: default, google, selector1, selector2, k1, mailgun, ses）",
		Fix: "請向您的郵件服務商取得 DKIM 設定，在 DNS 加入對應的 TXT 記錄（selector._domainkey." + domain + "）"}
}

func checkDMARC(domain string, score *int) ComplianceCheck {
	host := "_dmarc." + domain
	txts, err := net.LookupTXT(host)
	if err != nil {
		*score -= 20
		return ComplianceCheck{Name: "DMARC 記錄", Status: "fail", Detail: "無法查詢 _dmarc." + domain,
			Fix: fmt.Sprintf("在 DNS 加入 TXT 記錄 _dmarc.%s：v=DMARC1; p=quarantine; rua=mailto:dmarc@%s", domain, domain)}
	}
	for _, txt := range txts {
		if strings.HasPrefix(txt, "v=DMARC1") {
			status := "pass"
			if strings.Contains(txt, "p=none") {
				status = "warn"
				*score -= 10
				return ComplianceCheck{Name: "DMARC 記錄", Status: status, Detail: "已設定但策略為 none: " + txt,
					Fix: "建議將 p=none 改為 p=quarantine 或 p=reject 以提高信件可信度"}
			}
			return ComplianceCheck{Name: "DMARC 記錄", Status: status, Detail: "已設定: " + txt}
		}
	}
	*score -= 20
	return ComplianceCheck{Name: "DMARC 記錄", Status: "fail", Detail: "未找到 DMARC 記錄",
		Fix: fmt.Sprintf("在 DNS 加入 TXT 記錄 _dmarc.%s：v=DMARC1; p=quarantine; rua=mailto:dmarc@%s", domain, domain)}
}

func checkMX(domain string, score *int) ComplianceCheck {
	mxs, err := net.LookupMX(domain)
	if err != nil || len(mxs) == 0 {
		*score -= 10
		return ComplianceCheck{Name: "MX 記錄", Status: "warn", Detail: "未找到 MX 記錄，回信可能無法送達",
			Fix: "確認 " + domain + " 有設定 MX 記錄指向郵件伺服器"}
	}
	hosts := make([]string, len(mxs))
	for i, mx := range mxs {
		hosts[i] = mx.Host
	}
	return ComplianceCheck{Name: "MX 記錄", Status: "pass", Detail: "已設定: " + strings.Join(hosts, ", ")}
}

func checkReverseDNS(host string, score *int) ComplianceCheck {
	ips, err := net.LookupIP(host)
	if err != nil || len(ips) == 0 {
		return ComplianceCheck{Name: "反向 DNS (PTR)", Status: "warn", Detail: "無法解析 SMTP 主機 IP: " + host,
			Fix: "確認 SMTP 主機名稱可正確解析"}
	}
	ip := ips[0].String()
	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		*score -= 10
		return ComplianceCheck{Name: "反向 DNS (PTR)", Status: "warn", Detail: fmt.Sprintf("IP %s 無反向 DNS 記錄", ip),
			Fix: fmt.Sprintf("請向 ISP 或雲端供應商申請為 %s 設定 PTR 記錄指向 %s", ip, host)}
	}
	return ComplianceCheck{Name: "反向 DNS (PTR)", Status: "pass", Detail: fmt.Sprintf("IP %s → %s", ip, names[0])}
}
