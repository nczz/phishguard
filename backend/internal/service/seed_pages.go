package service

import "github.com/nczz/phishguard/internal/model"

func seedLandingPages(tenantID int64) []model.LandingPage {
	return []model.LandingPage{
		// 0: Microsoft 365 風格登入（密碼到期）
		{TenantID: &tenantID, Name: "密碼重設頁面", CaptureCredentials: true, CaptureFields: `["email","password"]`,
			HTML: `<!DOCTYPE html><html><head><meta charset="utf-8"><title>重設密碼</title>
<style>*{margin:0;padding:0;box-sizing:border-box}body{font-family:'Segoe UI',-apple-system,sans-serif;background:#f2f2f2;display:flex;justify-content:center;align-items:center;min-height:100vh}
.container{background:#fff;width:440px;padding:44px;box-shadow:0 2px 6px rgba(0,0,0,.2)}.logo{font-size:22px;font-weight:600;margin-bottom:24px;color:#1b1b1b}
.subtitle{color:#666;margin-bottom:20px;font-size:14px}input{width:100%;padding:10px 8px;margin:8px 0;border:none;border-bottom:2px solid #0078d4;font-size:15px;outline:none}
input:focus{border-bottom-color:#005a9e}button{width:100%;padding:12px;background:#0078d4;color:#fff;border:none;font-size:15px;cursor:pointer;margin-top:20px}
button:hover{background:#005a9e}.footer{margin-top:24px;font-size:12px;color:#999}</style></head>
<body><div class="container"><div class="logo">🔐 密碼重設</div>
<p class="subtitle">您的密碼即將到期，請輸入帳號資訊以重新設定密碼。</p>
<form action="{{.SubmitURL}}" method="POST">
<input name="email" placeholder="公司 Email" type="email" required>
<input name="password" placeholder="目前密碼" type="password" required>
<input name="new_password" placeholder="新密碼" type="password" required>
<button type="submit">確認重設</button></form>
<p class="footer">此頁面由 IT 服務台提供 | 如有疑問請撥分機 #2580</p></div></body></html>`},

		// 1: 包裹領取確認
		{TenantID: &tenantID, Name: "包裹領取確認", CaptureCredentials: true, CaptureFields: `["name","employee_id","phone"]`,
			HTML: `<!DOCTYPE html><html><head><meta charset="utf-8"><title>包裹領取確認</title>
<style>*{margin:0;padding:0;box-sizing:border-box}body{font-family:-apple-system,sans-serif;background:#f0f2f5;display:flex;justify-content:center;align-items:center;min-height:100vh}
.card{background:#fff;width:420px;border-radius:8px;box-shadow:0 2px 8px rgba(0,0,0,.1);overflow:hidden}
.header{background:#52c41a;color:#fff;padding:20px 24px;font-size:16px}.body{padding:24px}
.info{background:#f6ffed;border:1px solid #b7eb8f;border-radius:4px;padding:12px;margin-bottom:20px;font-size:13px}
input{width:100%;padding:10px;margin:6px 0;border:1px solid #d9d9d9;border-radius:4px;font-size:14px}
label{font-size:13px;color:#666;margin-top:8px;display:block}
button{width:100%;padding:12px;background:#52c41a;color:#fff;border:none;border-radius:4px;font-size:15px;cursor:pointer;margin-top:16px}</style></head>
<body><div class="card"><div class="header">📦 包裹領取確認</div><div class="body">
<div class="info">包裹編號：TW-2026-04 | 狀態：已送達 B1 收發室</div>
<p style="margin-bottom:12px;font-size:14px">請填寫以下資訊以確認領取：</p>
<form action="{{.SubmitURL}}" method="POST">
<label>姓名</label><input name="name" placeholder="請輸入姓名" required>
<label>員工編號</label><input name="employee_id" placeholder="如 EMP-001" required>
<label>聯絡電話</label><input name="phone" placeholder="分機或手機" required>
<button type="submit">確認領取</button></form></div></div></body></html>`},

		// 2: HR 薪資系統登入
		{TenantID: &tenantID, Name: "eHR 薪資系統登入", CaptureCredentials: true, CaptureFields: `["employee_id","password"]`,
			HTML: `<!DOCTYPE html><html><head><meta charset="utf-8"><title>eHR Portal - 登入</title>
<style>*{margin:0;padding:0;box-sizing:border-box}body{font-family:-apple-system,sans-serif;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);display:flex;justify-content:center;align-items:center;min-height:100vh}
.card{background:#fff;width:400px;border-radius:12px;box-shadow:0 20px 60px rgba(0,0,0,.3);overflow:hidden}
.header{background:linear-gradient(135deg,#722ed1,#9254de);color:#fff;padding:32px 24px;text-align:center}
.header h1{font-size:20px;margin-bottom:4px}.header p{font-size:13px;opacity:.8}.body{padding:32px 24px}
input{width:100%;padding:12px;margin:8px 0;border:1px solid #d9d9d9;border-radius:6px;font-size:14px}
button{width:100%;padding:12px;background:linear-gradient(135deg,#722ed1,#9254de);color:#fff;border:none;border-radius:6px;font-size:15px;cursor:pointer;margin-top:16px}
.note{text-align:center;margin-top:16px;font-size:12px;color:#999}</style></head>
<body><div class="card"><div class="header"><h1>eHR Portal</h1><p>人力資源管理系統</p></div><div class="body">
<form action="{{.SubmitURL}}" method="POST">
<input name="employee_id" placeholder="員工編號" required>
<input name="password" placeholder="密碼" type="password" required>
<button type="submit">登入</button></form>
<p class="note">忘記密碼？請聯繫人資部 分機 #1500</p></div></div></body></html>`},

		// 3: 安全驗證頁面（資安警告）
		{TenantID: &tenantID, Name: "帳號安全驗證", CaptureCredentials: true, CaptureFields: `["email","password"]`,
			HTML: `<!DOCTYPE html><html><head><meta charset="utf-8"><title>帳號安全驗證</title>
<style>*{margin:0;padding:0;box-sizing:border-box}body{font-family:-apple-system,sans-serif;background:#1a1a2e;display:flex;justify-content:center;align-items:center;min-height:100vh}
.card{background:#fff;width:440px;border-radius:8px;box-shadow:0 20px 60px rgba(0,0,0,.5);overflow:hidden}
.alert{background:#ff4d4f;color:#fff;padding:20px 24px;text-align:center}.alert h2{font-size:18px;margin-bottom:4px}.alert p{font-size:13px;opacity:.9}
.body{padding:28px 24px}
.warning{background:#fff2f0;border:1px solid #ffccc7;border-radius:4px;padding:12px;margin-bottom:20px;font-size:13px;color:#cf1322}
input{width:100%;padding:12px;margin:8px 0;border:1px solid #d9d9d9;border-radius:4px;font-size:14px}
button{width:100%;padding:12px;background:#ff4d4f;color:#fff;border:none;border-radius:4px;font-size:15px;cursor:pointer;margin-top:16px}
.timer{text-align:center;color:#ff4d4f;font-size:13px;margin-top:12px}</style></head>
<body><div class="card"><div class="alert"><h2>🔒 帳號安全驗證</h2><p>偵測到異常登入，請驗證您的身份</p></div>
<div class="body">
<div class="warning">⚠️ 您的帳號於陌生裝置（越南 胡志明市）登入，若非本人操作請立即驗證。</div>
<form action="{{.SubmitURL}}" method="POST">
<input name="email" placeholder="公司 Email" type="email" required>
<input name="password" placeholder="密碼" type="password" required>
<button type="submit">驗證身份</button></form>
<p class="timer">⏱ 此驗證頁面將於 15 分鐘後失效</p></div></div></body></html>`},

		// 4: 財務系統登入（發票簽核）
		{TenantID: &tenantID, Name: "財務系統登入", CaptureCredentials: true, CaptureFields: `["email","password"]`,
			HTML: `<!DOCTYPE html><html><head><meta charset="utf-8"><title>財務管理系統 - 登入</title>
<style>*{margin:0;padding:0;box-sizing:border-box}body{font-family:-apple-system,sans-serif;background:#f5f5f5;display:flex;justify-content:center;align-items:center;min-height:100vh}
.card{background:#fff;width:420px;border-radius:8px;box-shadow:0 2px 12px rgba(0,0,0,.1);overflow:hidden}
.header{background:linear-gradient(135deg,#fa8c16,#faad14);color:#fff;padding:24px;text-align:center}
.header h1{font-size:18px;margin-bottom:4px}.header p{font-size:12px;opacity:.9}.body{padding:24px}
.invoice{background:#fff7e6;border:1px solid #ffd591;border-radius:4px;padding:12px;margin-bottom:20px;font-size:13px}
input{width:100%;padding:12px;margin:8px 0;border:1px solid #d9d9d9;border-radius:4px;font-size:14px}
button{width:100%;padding:12px;background:#fa8c16;color:#fff;border:none;border-radius:4px;font-size:15px;cursor:pointer;margin-top:16px}
.note{text-align:center;margin-top:16px;font-size:12px;color:#999}</style></head>
<body><div class="card"><div class="header"><h1>💰 財務管理系統</h1><p>Finance Management System</p></div><div class="body">
<div class="invoice">📄 待簽核發票：AB-20260401-0037 | 金額：NT$ 28,500 | 到期：2026/04/20</div>
<p style="font-size:14px;margin-bottom:12px">請登入以查看發票詳情並完成簽核：</p>
<form action="{{.SubmitURL}}" method="POST">
<input name="email" placeholder="公司 Email" type="email" required>
<input name="password" placeholder="密碼" type="password" required>
<button type="submit">登入查看</button></form>
<p class="note">財務部 | 分機 #3100</p></div></div></body></html>`},
	}
}
