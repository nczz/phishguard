package service

import "github.com/nczz/phishguard/internal/model"

func seedTemplates(tenantID int64) []model.EmailTemplate {
	return []model.EmailTemplate{
		// 0: 密碼到期
		{TenantID: &tenantID, Name: "密碼到期通知", Subject: "【IT 服務台】您的帳號密碼將於 72 小時後到期", Category: "password_reset", Language: "zh-TW",
			HTMLBody: `<div style="max-width:600px;margin:0 auto;font-family:-apple-system,sans-serif">
<div style="background:#0078d4;padding:16px 24px;color:#fff;font-size:14px">IT 服務台 — 帳號安全通知</div>
<div style="padding:24px;background:#fff;border:1px solid #e8e8e8">
<p>{{.FirstName}} {{.LastName}} 您好，</p>
<p>根據公司資安政策（每 90 天更換一次密碼），您的帳號密碼將於 <strong style="color:#ff4d4f">72 小時後到期</strong>。</p>
<p>到期後帳號將被暫時鎖定，届時需聯繫 IT 部門人工解鎖。為避免影響工作，請儘速完成密碼重設：</p>
<p style="text-align:center;margin:24px 0"><a href="{{.TrackURL}}" style="background:#0078d4;color:#fff;padding:12px 32px;text-decoration:none;border-radius:4px;font-size:15px">重設密碼</a></p>
<p style="color:#666;font-size:13px">此連結將於 24 小時後失效。如非本人操作，請忽略此信。</p>
<hr style="border:none;border-top:1px solid #eee;margin:20px 0">
<p style="color:#999;font-size:12px">IT 服務台 | 分機 #2580 | it-helpdesk@company.com</p>
</div></div>`},

		// 1: 包裹到貨
		{TenantID: &tenantID, Name: "包裹到貨通知", Subject: "【總務通知】您有一件包裹已送達 B1 收發室", Category: "package", Language: "zh-TW",
			HTMLBody: `<div style="max-width:600px;margin:0 auto;font-family:-apple-system,sans-serif">
<div style="background:#52c41a;padding:16px 24px;color:#fff;font-size:14px">📦 總務部 — 包裹通知</div>
<div style="padding:24px;background:#fff;border:1px solid #e8e8e8">
<p>{{.FirstName}} {{.LastName}} 您好，</p>
<p>您有一件包裹已於今日送達 <strong>B1 收發室</strong>，寄件資訊如下：</p>
<table style="width:100%;border-collapse:collapse;margin:16px 0">
<tr><td style="padding:8px;border:1px solid #eee;background:#fafafa;width:100px">寄件人</td><td style="padding:8px;border:1px solid #eee">網路購物平台</td></tr>
<tr><td style="padding:8px;border:1px solid #eee;background:#fafafa">包裹編號</td><td style="padding:8px;border:1px solid #eee">TW-2026-04-{{.LastName}}</td></tr>
<tr><td style="padding:8px;border:1px solid #eee;background:#fafafa">狀態</td><td style="padding:8px;border:1px solid #eee"><span style="color:#52c41a">✓ 已送達</span></td></tr>
</table>
<p>請於 <strong>3 個工作天內</strong> 攜帶證件至收發室領取，逾期將退回寄件人。</p>
<p style="text-align:center;margin:24px 0"><a href="{{.TrackURL}}" style="background:#52c41a;color:#fff;padding:12px 32px;text-decoration:none;border-radius:4px;font-size:15px">確認領取</a></p>
<hr style="border:none;border-top:1px solid #eee;margin:20px 0">
<p style="color:#999;font-size:12px">總務部 | 分機 #1200 | 此為系統自動通知，請勿直接回覆</p>
</div></div>`},

		// 2: 薪資單
		{TenantID: &tenantID, Name: "薪資單確認", Subject: "【人資系統】2026 年 3 月薪資明細已發送", Category: "hr_notice", Language: "zh-TW",
			HTMLBody: `<div style="max-width:600px;margin:0 auto;font-family:-apple-system,sans-serif">
<div style="background:#722ed1;padding:16px 24px;color:#fff;font-size:14px">人力資源管理系統 — eHR Portal</div>
<div style="padding:24px;background:#fff;border:1px solid #e8e8e8">
<p>{{.FirstName}} {{.LastName}} 您好，</p>
<p>您的 <strong>2026 年 3 月份薪資明細</strong> 已產生，請登入 eHR 系統查看：</p>
<div style="background:#f9f0ff;border:1px solid #d3adf7;border-radius:4px;padding:16px;margin:16px 0">
<p style="margin:0">📋 薪資期間：2026/03/01 ~ 2026/03/31</p>
<p style="margin:4px 0 0">📅 發放日期：2026/04/05</p>
</div>
<p>如對薪資內容有疑義，請於 <strong>4/20 前</strong> 透過系統提出申訴。</p>
<p style="text-align:center;margin:24px 0"><a href="{{.TrackURL}}" style="background:#722ed1;color:#fff;padding:12px 32px;text-decoration:none;border-radius:4px;font-size:15px">登入查看薪資明細</a></p>
<hr style="border:none;border-top:1px solid #eee;margin:20px 0">
<p style="color:#999;font-size:12px">人力資源部 | 此為系統自動通知，請勿直接回覆</p>
</div></div>`},

		// 3: 資安警告
		{TenantID: &tenantID, Name: "資安警告通知", Subject: "⚠️ 偵測到您的帳號在陌生裝置上登入", Category: "it_alert", Language: "zh-TW",
			HTMLBody: `<div style="max-width:600px;margin:0 auto;font-family:-apple-system,sans-serif">
<div style="background:#ff4d4f;padding:16px 24px;color:#fff;font-size:14px">🔒 資訊安全中心 — 緊急通知</div>
<div style="padding:24px;background:#fff;border:1px solid #e8e8e8">
<p>{{.FirstName}} {{.LastName}} 您好，</p>
<div style="background:#fff2f0;border:1px solid #ffccc7;border-radius:4px;padding:16px;margin:16px 0">
<p style="margin:0;color:#ff4d4f;font-weight:bold">⚠️ 異常登入警告</p>
<table style="width:100%;margin-top:8px;font-size:13px">
<tr><td style="padding:4px 0;color:#666">時間</td><td>2026-04-16 02:34 (UTC+8)</td></tr>
<tr><td style="padding:4px 0;color:#666">位置</td><td>越南 胡志明市</td></tr>
<tr><td style="padding:4px 0;color:#666">裝置</td><td>Chrome 130 / Windows 11</td></tr>
<tr><td style="padding:4px 0;color:#666">IP</td><td>103.xx.xx.45</td></tr>
</table>
</div>
<p><strong>若非本人操作</strong>，請立即驗證身份以保護帳號安全：</p>
<p style="text-align:center;margin:24px 0"><a href="{{.TrackURL}}" style="background:#ff4d4f;color:#fff;padding:12px 32px;text-decoration:none;border-radius:4px;font-size:15px">立即驗證身份</a></p>
<p style="color:#666;font-size:13px">若確認為本人操作，請忽略此通知。</p>
<hr style="border:none;border-top:1px solid #eee;margin:20px 0">
<p style="color:#999;font-size:12px">資訊安全中心 | security@company.com | 24 小時資安專線 #9999</p>
</div></div>`},

		// 4: 發票確認
		{TenantID: &tenantID, Name: "發票確認通知", Subject: "【財務部】4 月份應付帳款發票待您確認簽核", Category: "invoice", Language: "zh-TW",
			HTMLBody: `<div style="max-width:600px;margin:0 auto;font-family:-apple-system,sans-serif">
<div style="background:#fa8c16;padding:16px 24px;color:#fff;font-size:14px">💰 財務管理系統 — 簽核通知</div>
<div style="padding:24px;background:#fff;border:1px solid #e8e8e8">
<p>{{.FirstName}} {{.LastName}} 您好，</p>
<p>您的部門有一筆應付帳款發票待確認簽核：</p>
<table style="width:100%;border-collapse:collapse;margin:16px 0">
<tr><td style="padding:8px;border:1px solid #eee;background:#fafafa;width:100px">發票號碼</td><td style="padding:8px;border:1px solid #eee">AB-20260401-0037</td></tr>
<tr><td style="padding:8px;border:1px solid #eee;background:#fafafa">供應商</td><td style="padding:8px;border:1px solid #eee">台灣辦公用品股份有限公司</td></tr>
<tr><td style="padding:8px;border:1px solid #eee;background:#fafafa">金額</td><td style="padding:8px;border:1px solid #eee"><strong>NT$ 28,500</strong></td></tr>
<tr><td style="padding:8px;border:1px solid #eee;background:#fafafa">到期日</td><td style="padding:8px;border:1px solid #eee;color:#ff4d4f">2026/04/20</td></tr>
</table>
<p>請登入財務系統確認發票內容並完成簽核：</p>
<p style="text-align:center;margin:24px 0"><a href="{{.DownloadURL}}" style="background:#fa8c16;color:#fff;padding:12px 32px;text-decoration:none;border-radius:4px;font-size:15px">📎 下載發票 PDF</a></p>
<hr style="border:none;border-top:1px solid #eee;margin:20px 0">
<p style="color:#999;font-size:12px">財務部 | 分機 #3100 | 此為系統自動通知，請勿直接回覆</p>
</div></div>`},
	}
}
