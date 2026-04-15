# 系統架構

## 總覽

PhishGuard 由三個 Go binary + React SPA + Nginx 組成。

```
                    ┌──────────────────────────────────┐
                    │        Nginx (80/443)            │
                    │  app.domain → Frontend + API     │
                    │  t.domain   → Tracker            │
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
             │  (mail +   │
             │  scheduler)│
             └────────────┘
```

## 元件職責

### API Server (`cmd/api`)
- 平台管理員 API (`/api/admin/*`) — 租戶管理、使用者管理、跨租戶稽核、方案設定
- 租戶 API (`/api/*`) — Campaign CRUD、模板、情境、收件人、報表、SMTP、合規檢測
- 認證（JWT 24h）、租戶隔離 middleware、稽核日誌 middleware

### Track Server (`cmd/tracker`)
- `GET /t/o/:rid` — 開信追蹤（1x1 透明 GIF）
- `GET /t/c/:rid` — 點擊追蹤（302 redirect 到 Landing Page）
- `GET /t/d/:rid/:filename` — 附件下載追蹤
- `POST /t/s/:rid` — 表單提交追蹤（顯示教育頁）
- `GET /t/r/:rid` — 舉報追蹤（顯示感謝頁）
- `GET /t/landing?rid=` — Landing Page 渲染

### Worker (`cmd/worker`)
- 每 30 秒輪詢發信佇列（撈 send_date <= now 的 results）
- 按 provider 政策限速（SES 12/sec, Mailgun 40/sec, SMTP 3/sec）
- 同域名每小時上限（防止 flooding）
- 暫時性錯誤自動重試（421/450/throttle）
- Campaign 完成後自動寄送報表給租戶管理員
- 定期測試排程器（每 30 秒檢查 auto_test_configs）

### Frontend (React SPA)
- 平台管理員 Dashboard（租戶統計、使用者管理、方案設定、切換租戶視角）
- 租戶介面（三步驟精靈、報表、收件人管理、情境庫、設定）
- 使用指南 + 流程總覽

## 發信架構

```
Campaign Launch
    ↓
選人引擎（冷卻期排除 → 全選/部門/隨機抽樣）
    ↓
排程引擎（立即/排程，工作時間/避開週末，均勻分配+抖動）
    ↓
Worker 輪詢（每 30 秒）
    ↓
Rate Limiter（按 provider 政策）
├── SES:     12/sec, 80ms interval, burst 14
├── Mailgun: 40/sec, 25ms interval, burst 50
└── SMTP:    3/sec, 350ms interval, burst 5
    ↓
同域名限制（SES 500/hr, Mailgun 500/hr, SMTP 100/hr）
    ↓
發送（含 List-Unsubscribe, Message-ID, 追蹤 pixel, 舉報連結）
    ↓
暫時性錯誤 → 保持 scheduled，下輪重試
永久性錯誤 → 標記 error + 記錄原因
```

## 多租戶隔離

- Application Layer：所有 repo 方法強制帶 tenant_id
- Middleware：從 JWT 取 tenant_id 注入 context
- 平台管理員：tenant_id=NULL，可透過 X-Tenant-ID header 或 impersonate 切換

## 方案限制

| | Free | Pro | Enterprise |
|---|:---:|:---:|:---:|
| 收件人 | 50 | 1,000 | 無限 |
| 年度活動 | 4 | 無限 | 無限 |
| 月發信量 | 200 | 10,000 | 無限 |
| 自訂模板 | ❌ | ✅ | ✅ |
| 自動測試 | ❌ | ✅ | ✅ |
| 部門報表 | ❌ | ✅ | ✅ |

管理員可在租戶詳情頁覆蓋任何限制值。

## 部署架構

```
單機 VPS (2GB RAM)
├── Nginx (reverse proxy + TLS)
├── API Server (Go binary)
├── Tracker Server (Go binary)
├── Worker (Go binary)
├── MySQL 8.0 (persistent volume)
└── Redis 7 (persistent volume)

DNS:
  app.yourdomain.com → VPS IP
  t.yourdomain.com   → VPS IP
```
