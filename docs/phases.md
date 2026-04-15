# 階段規劃與驗收條件

## Phase 1 — 核心引擎 ✅

- [x] MySQL schema（16 張表）
- [x] 租戶 CRUD + 建立時自動 seed 範例資料
- [x] 使用者 CRUD + JWT 認證
- [x] 情境庫 CRUD（模板 + Landing Page + 教育頁）
- [x] 收件人 CSV 匯入（upsert + 性別欄位 + 批次刪除）
- [x] SMTP Profile CRUD + 測試發信 + 合規檢測（SPF/DKIM/DMARC）
- [x] Campaign 建立 + 選人引擎（冷卻期 → 全選/隨機/部門）
- [x] 排程引擎（立即/排程、工作時間、避開週末、均勻分配+抖動）
- [x] Mail Worker（SMTP + Mailgun + SES、provider 限速、暫時性錯誤重試）
- [x] Track Server（開信/點擊/下載/提交/舉報 + 教育頁 + 感謝頁）
- [x] 漏斗統計 + 部門統計 API
- [x] 稽核日誌 middleware

## Phase 2 — 前端 UI ✅

- [x] Login（角色自動導向）
- [x] 三步驟 Campaign Wizard（選情境 → 選對象+排程 → 確認發送）
- [x] Campaign Detail（進度條 + 漏斗圖 + 部門排名 + 收件人明細）
- [x] Tenant Dashboard（真實統計：平均開信/點擊/提交率）
- [x] 收件人管理（CSV 匯入 + 部門統計 + 編輯/刪除/批次刪除 + 匯出）
- [x] 情境庫（卡片瀏覽 + CRUD + 教育頁預覽）
- [x] 模板管理（CRUD + Drawer 編輯器 + 欄位 tooltip）
- [x] Landing Page 管理（CRUD + HTML 預覽）
- [x] SMTP 設定（多 mailer 類型 + 測試發信 + 合規檢測面板）
- [x] 稽核日誌（分頁 + 展開詳情）
- [x] 使用指南（5 步驟 + 變數表 + 指標說明 + FAQ）
- [x] 流程總覽（視覺化泳道圖）
- [x] 匯入範例資料頁面
- [x] 自動測試設定 UI
- [x] 平台管理員 Dashboard（5 統計 + 警示 + 租戶表）
- [x] 租戶詳情（啟停用 + 切換視角 + 活動列表 + 使用者管理 + 方案設定）
- [x] 租戶列表 + 建立租戶
- [x] 跨租戶稽核日誌
- [x] 租戶切換（impersonate + banner + 返回）

## Phase 3 — 報表與自動化 ✅

- [x] Dashboard 真實統計（平均開信/點擊/提交率）
- [x] Campaign 收件人明細表（每人狀態 + 時間戳）
- [x] CSV 結果匯出
- [x] PDF 報表匯出（HTML 報表 + 瀏覽器列印）
- [x] 報表自動寄送（Campaign 完成後 email）
- [x] 累犯追蹤（跨 Campaign 個人歷史 + 風險等級）
- [x] 趨勢分析（折線圖 4 指標）
- [x] 自動定期測試排程器（月/季/半年）
- [x] 冷卻期（30 天內不重複，排除先於抽樣）
- [x] 方案限制系統（收件人/活動/發信量 + 管理員覆蓋）
- [x] 發信合規（SPF/DKIM/DMARC 檢測 + provider 限速 + 同域名限制）
- [x] 舉報連結自動注入 + 感謝頁
- [x] 發送進度條（即時 + 預估剩餘時間）
- [x] API 回應格式規範 + 統一 helper

## Phase 4 — 企業級功能（未實作）

- [ ] SSO/OIDC（Azure AD, Okta）
- [ ] 自訂發信域名（per-tenant wildcard TLS）
- [ ] 白標（自訂 logo/色彩）
- [ ] API Key 認證（供客戶程式化操作）
- [ ] Outlook/Gmail 舉報外掛
- [ ] 產業對比數據（跨租戶匿名統計）

## 部署

- [x] Docker Compose 開發環境
- [x] Docker Compose 生產環境（6 容器）
- [x] Dockerfile（Go multi-stage + Frontend multi-stage）
- [x] Nginx reverse proxy 設定
- [x] VPS 一鍵部署腳本（自動產生密碼）
