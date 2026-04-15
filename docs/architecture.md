# 系統架構

## 總覽

PhishGuard 由三個獨立的 Go binary + 一個 React SPA 組成，共用同一個 MySQL 資料庫和 Redis。

```
                        ┌──────────────────────────────────┐
                        │          Load Balancer           │
                        └──────┬──────────────┬────────────┘
                               │              │
                    ┌──────────▼───┐   ┌──────▼──────────┐
                    │   API Server │   │  Track Server   │
                    │   (Gin)      │   │  (net/http)     │
                    │   :8080      │   │  :8090          │
                    └──────┬───┬──┘   └───┬──────────────┘
                           │   │          │
                    ┌──────▼┐ ┌▼──────────▼──┐
                    │ Redis │ │    MySQL      │
                    │ :6379 │ │    :3306      │
                    └──┬────┘ └──────────────┘
                       │
                 ┌─────▼──────┐
                 │   Worker   │
                 │  (asynq)   │
                 └────────────┘
```

## 元件職責

### API Server (`cmd/api`)

- 平台管理員 API (`/api/admin/*`)
- 租戶 API (`/api/*`)
- 認證（JWT）
- 靜態檔案服務（React build）

### Track Server (`cmd/tracker`)

- 開信追蹤 `GET /t/o/:rid` — 回傳 1x1 透明 GIF
- 點擊追蹤 `GET /t/c/:rid` — 記錄後 302 redirect
- 附件下載追蹤 `GET /t/d/:rid/:filename` — 記錄後回傳檔案
- 表單提交追蹤 `POST /t/s/:rid` — 記錄後顯示教育頁
- 舉報回報 `POST /t/r/:rid` — 記錄舉報事件

設計原則：極輕量、無認證、低延遲。與 API Server 分離，避免互相影響。

### Worker (`cmd/worker`)

- 消費 Redis queue 發送郵件（SMTP / Mailgun / SES）
- 定期測試排程檢查（每分鐘）
- Campaign 完成後自動產生報表 + 寄送

### Frontend (React SPA)

- 平台管理員 Dashboard
- 租戶管理介面（Campaign、情境、收件人、報表、稽核日誌）

## 資料流

### Campaign 執行流程

```
客戶建立 Campaign
       │
       ▼
API Server: 驗證 → 選人引擎（抽樣/部門篩選/去重）
       │
       ▼
為每個收件人產生 rid (UUID) → 寫入 results 表
       │
       ▼
渲染信件（模板變數替換、嵌入追蹤 pixel、改寫連結）
       │
       ▼
推入 Redis queue（每封信一個 task）
       │
       ▼
Worker 消費 → 透過 mailer 抽象層發信
       │
       ├── SMTP → 直連 SMTP server
       ├── Mailgun → Mailgun API
       └── SES → AWS SES API
       │
       ▼
更新 results.status = sent / error
```

### 追蹤流程

```
收件人開信 → 載入 tracking pixel
       │
       ▼
Track Server: GET /t/o/:rid
       │
       ▼
查 results 表取得 rid → 寫入 events 表 → 更新 results.opened_at
       │
收件人點連結 → GET /t/c/:rid
       │
       ▼
記錄 clicked → 302 redirect 到 Landing Page（帶 rid）
       │
收件人提交表單 → POST /t/s/:rid
       │
       ▼
記錄 submitted（只記欄位名，不記值）→ 回傳教育頁 HTML
```

## 多租戶隔離

### Application Layer（主要防線）

所有 DB 存取經過 `repo` 層，強制帶 `tenant_id`：

```go
// repo 層的每個方法都要求 tenantID 參數
func (r *CampaignRepo) FindAll(ctx context.Context, tenantID int64) ([]model.Campaign, error)
```

### Middleware Layer

Tenant middleware 從 JWT 中取出 tenant_id，注入 context。所有 handler 從 context 取得，不接受客戶端傳入。

### 平台管理員

平台管理員的 tenant_id 為 NULL，可透過 `X-Tenant-ID` header 切換到任意租戶操作。

## 認證架構

```
登入 → POST /api/auth/login → 驗證帳密 → 回傳 JWT
                                          │
                                          ▼
                                JWT payload:
                                {
                                  "sub": user_id,
                                  "tid": tenant_id,  // NULL for platform admin
                                  "role": "tenant_admin",
                                  "exp": ...
                                }
                                          │
                                          ▼
                        後續請求帶 Authorization: Bearer <token>
                                          │
                                          ▼
                        Auth middleware 驗證 → Tenant middleware 注入 tenant_id
```

## 發信架構（Mailer 抽象層）

```go
// Mailer interface — 所有發信方式實作此介面
type Mailer interface {
    Send(ctx context.Context, msg *Message) error
    Name() string
}

// 三種實作
├── SMTPMailer      // 直連 SMTP server
├── MailgunMailer   // Mailgun HTTP API
└── SESMailer       // AWS SES SDK
```

租戶的 `smtp_profiles` 表記錄使用哪種 mailer 及對應的認證資訊。Worker 發信時根據 profile 選擇對應的 Mailer 實作。

## 部署架構（Production）

```
┌─ Docker Compose / K8s ─────────────────────────────┐
│                                                     │
│  nginx (reverse proxy + TLS)                        │
│  ├── app.phishguard.tw    → api:8080                │
│  └── t.phishguard.tw      → tracker:8090            │
│                                                     │
│  api      ×2 replicas                               │
│  tracker  ×2 replicas                               │
│  worker   ×1~N (依發信量 auto-scale)                 │
│                                                     │
│  MySQL 8.0 (persistent volume)                      │
│  Redis 7   (persistent volume)                      │
└─────────────────────────────────────────────────────┘
```
