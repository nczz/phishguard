package service

func seedEducationHTML() string {
	return `<!DOCTYPE html><html><head><meta charset="utf-8"><title>釣魚測試結果</title>
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
}
