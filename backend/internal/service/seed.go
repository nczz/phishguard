package service

import (
	"github.com/phishguard/phishguard/internal/model"
	"gorm.io/gorm"
)

func SeedTenantData(db *gorm.DB, tenantID int64) error {
	// --- Email Templates ---
	templates := []model.EmailTemplate{
		{TenantID: &tenantID, Name: "密碼到期通知", Subject: "【重要】您的密碼即將到期", Category: "password_reset", Language: "zh-TW",
			HTMLBody: `<p>親愛的 {{.FirstName}} {{.LastName}} 您好，</p>
<p>您的公司帳號密碼將於 <strong>3 天後到期</strong>，請立即點擊下方連結重新設定密碼，以免帳號被鎖定。</p>
<p><a href="{{.TrackURL}}" style="background:#1677ff;color:#fff;padding:10px 24px;text-decoration:none;border-radius:4px;">立即重設密碼</a></p>
<p>如有疑問請聯繫 IT 部門。</p><p>IT 服務台</p>`},
		{TenantID: &tenantID, Name: "包裹到貨通知", Subject: "您有一個包裹待領取", Category: "package", Language: "zh-TW",
			HTMLBody: `<p>{{.FirstName}} {{.LastName}} 您好，</p>
<p>您有一個包裹已送達公司大廳，請點擊下方連結確認領取資訊：</p>
<p><a href="{{.TrackURL}}" style="background:#52c41a;color:#fff;padding:10px 24px;text-decoration:none;border-radius:4px;">確認包裹資訊</a></p>
<p>若非本人包裹請忽略此信。</p><p>總務部</p>`},
		{TenantID: &tenantID, Name: "薪資單確認", Subject: "本月薪資明細已發送，請確認", Category: "hr_notice", Language: "zh-TW",
			HTMLBody: `<p>{{.FirstName}} {{.LastName}} 您好，</p>
<p>您的本月薪資明細已產生，請點擊下方連結登入查看並確認：</p>
<p><a href="{{.TrackURL}}" style="background:#faad14;color:#fff;padding:10px 24px;text-decoration:none;border-radius:4px;">查看薪資明細</a></p>
<p>如有異議請於 3 個工作天內回覆。</p><p>人力資源部</p>`},
		{TenantID: &tenantID, Name: "資安警告通知", Subject: "偵測到異常登入活動", Category: "it_alert", Language: "zh-TW",
			HTMLBody: `<p>{{.FirstName}} {{.LastName}} 您好，</p>
<p>系統偵測到您的帳號於 <strong>未知裝置</strong> 上有異常登入嘗試。若非本人操作，請立即點擊下方連結驗證身份：</p>
<p><a href="{{.TrackURL}}" style="background:#ff4d4f;color:#fff;padding:10px 24px;text-decoration:none;border-radius:4px;">立即驗證</a></p>
<p>資訊安全中心</p>`},
		{TenantID: &tenantID, Name: "發票確認通知", Subject: "請確認附件發票內容", Category: "invoice", Language: "zh-TW",
			HTMLBody: `<p>{{.FirstName}} {{.LastName}} 您好，</p>
<p>附件為本月應付發票，請點擊下方連結下載並確認金額：</p>
<p><a href="{{.TrackURL}}" style="background:#722ed1;color:#fff;padding:10px 24px;text-decoration:none;border-radius:4px;">📎 下載發票 (PDF)</a></p>
<p>請於本週五前完成確認。</p><p>財務部</p>`},
	}
	for i := range templates {
		if err := db.Create(&templates[i]).Error; err != nil {
			return err
		}
	}

	// --- Landing Pages ---
	loginPage := model.LandingPage{
		TenantID: &tenantID, Name: "仿登入頁面", CaptureCredentials: true,
		CaptureFields: `["email","password"]`,
		HTML: `<!DOCTYPE html><html><head><meta charset="utf-8"><title>帳號登入</title>
<style>body{font-family:-apple-system,sans-serif;display:flex;justify-content:center;align-items:center;min-height:100vh;background:#f0f2f5;margin:0}
.card{background:#fff;padding:40px;border-radius:8px;box-shadow:0 2px 8px rgba(0,0,0,.1);width:360px}
h2{text-align:center;margin-bottom:24px}input{width:100%;padding:10px;margin:8px 0;border:1px solid #d9d9d9;border-radius:4px;box-sizing:border-box}
button{width:100%;padding:10px;background:#1677ff;color:#fff;border:none;border-radius:4px;cursor:pointer;font-size:16px;margin-top:12px}</style></head>
<body><div class="card"><h2>🔐 帳號登入</h2><form action="{{.SubmitURL}}" method="POST">
<input name="email" placeholder="Email" type="email" required>
<input name="password" placeholder="密碼" type="password" required>
<button type="submit">登入</button></form></div></body></html>`,
	}
	if err := db.Create(&loginPage).Error; err != nil {
		return err
	}

	confirmPage := model.LandingPage{
		TenantID: &tenantID, Name: "確認資訊頁面", CaptureCredentials: true,
		CaptureFields: `["name","employee_id"]`,
		HTML: `<!DOCTYPE html><html><head><meta charset="utf-8"><title>身份確認</title>
<style>body{font-family:-apple-system,sans-serif;display:flex;justify-content:center;align-items:center;min-height:100vh;background:#f0f2f5;margin:0}
.card{background:#fff;padding:40px;border-radius:8px;box-shadow:0 2px 8px rgba(0,0,0,.1);width:360px}
h2{text-align:center;margin-bottom:24px}input{width:100%;padding:10px;margin:8px 0;border:1px solid #d9d9d9;border-radius:4px;box-sizing:border-box}
button{width:100%;padding:10px;background:#52c41a;color:#fff;border:none;border-radius:4px;cursor:pointer;font-size:16px;margin-top:12px}</style></head>
<body><div class="card"><h2>📋 身份確認</h2><form action="{{.SubmitURL}}" method="POST">
<input name="name" placeholder="姓名" required>
<input name="employee_id" placeholder="員工編號" required>
<button type="submit">確認</button></form></div></body></html>`,
	}
	if err := db.Create(&confirmPage).Error; err != nil {
		return err
	}

	// --- Education HTML ---
	educationHTML := `<!DOCTYPE html><html><head><meta charset="utf-8"><title>釣魚測試結果</title>
<style>body{font-family:-apple-system,sans-serif;display:flex;justify-content:center;align-items:center;min-height:100vh;background:#f0f2f5;margin:0}
.card{background:#fff;padding:40px;border-radius:8px;box-shadow:0 2px 8px rgba(0,0,0,.1);max-width:600px}
h1{color:#ff4d4f;text-align:center}h2{color:#1677ff}.tip{background:#fff7e6;border-left:4px solid #faad14;padding:12px;margin:12px 0;border-radius:4px}</style></head>
<body><div class="card">
<h1>⚠️ 這是一封釣魚測試信</h1>
<p>您剛才點擊的連結是公司資安團隊發送的<strong>釣魚模擬測試</strong>，目的是提升全體員工的資安意識。</p>
<h2>如何辨識釣魚信？</h2>
<div class="tip">🔍 <strong>檢查寄件者</strong>：注意 email 地址是否為公司官方域名</div>
<div class="tip">🔗 <strong>檢查連結</strong>：滑鼠移到連結上方，確認網址是否正確</div>
<div class="tip">⏰ <strong>注意緊迫感</strong>：釣魚信常用「立即」「到期」等字眼製造壓力</div>
<div class="tip">📎 <strong>小心附件</strong>：不要開啟來路不明的附件</div>
<div class="tip">🚨 <strong>遇到可疑信件</strong>：請使用郵件舉報按鈕或通知 IT 部門</div>
<p style="text-align:center;margin-top:24px;color:#999">本測試由 PhishGuard 釣魚模擬平台提供</p>
</div></body></html>`

	// --- Scenarios ---
	scenarios := []model.Scenario{
		{TenantID: &tenantID, Name: "密碼到期通知", Category: "password_reset", Difficulty: 2, Language: "zh-TW",
			TemplateID: &templates[0].ID, PageID: &loginPage.ID, EducationHTML: educationHTML, IsActive: true},
		{TenantID: &tenantID, Name: "包裹到貨通知", Category: "package", Difficulty: 1, Language: "zh-TW",
			TemplateID: &templates[1].ID, PageID: &confirmPage.ID, EducationHTML: educationHTML, IsActive: true},
		{TenantID: &tenantID, Name: "薪資單確認", Category: "hr_notice", Difficulty: 2, Language: "zh-TW",
			TemplateID: &templates[2].ID, PageID: &loginPage.ID, EducationHTML: educationHTML, IsActive: true},
		{TenantID: &tenantID, Name: "資安警告通知", Category: "it_alert", Difficulty: 3, Language: "zh-TW",
			TemplateID: &templates[3].ID, PageID: &loginPage.ID, EducationHTML: educationHTML, IsActive: true},
		{TenantID: &tenantID, Name: "發票確認通知", Category: "invoice", Difficulty: 2, Language: "zh-TW",
			TemplateID: &templates[4].ID, PageID: &confirmPage.ID, EducationHTML: educationHTML, IsActive: true},
	}
	for i := range scenarios {
		if err := db.Create(&scenarios[i]).Error; err != nil {
			return err
		}
	}

	// --- Sample Recipient Group ---
	group := model.RecipientGroup{TenantID: tenantID, Name: "範例員工"}
	if err := db.Create(&group).Error; err != nil {
		return err
	}
	sampleRecipients := []model.Recipient{
		{TenantID: tenantID, GroupID: group.ID, Email: "demo-wang@example.com", FirstName: "小明", LastName: "王", Department: "業務部", Gender: "男", Position: "業務經理"},
		{TenantID: tenantID, GroupID: group.ID, Email: "demo-chen@example.com", FirstName: "小華", LastName: "陳", Department: "財務部", Gender: "女", Position: "會計"},
		{TenantID: tenantID, GroupID: group.ID, Email: "demo-lin@example.com", FirstName: "小美", LastName: "林", Department: "研發部", Gender: "女", Position: "工程師"},
		{TenantID: tenantID, GroupID: group.ID, Email: "demo-zhang@example.com", FirstName: "大偉", LastName: "張", Department: "行政部", Gender: "男", Position: "行政專員"},
		{TenantID: tenantID, GroupID: group.ID, Email: "demo-li@example.com", FirstName: "雅婷", LastName: "李", Department: "業務部", Gender: "不指定", Position: "業務代表"},
	}
	return db.CreateInBatches(sampleRecipients, 100).Error
}
