# PhishGuard

多租戶釣魚模擬測試 SaaS 平台。協助企業定期執行釣魚演練、量化員工資安風險、產出稽核報告。

## 核心功能

- **多租戶架構** — 平台管理員為企業客戶開設獨立帳號，資料完全隔離
- **方案管理** — Free/Pro/Enterprise 三級方案，收件人/活動/發信量限制，管理員可覆蓋
- **情境庫** — 預建釣魚情境（信件模板 + Landing Page + 教育頁），租戶建立時自動 seed
- **自訂模板** — HTML 信件編輯器 + Landing Page 編輯器 + 教育頁預覽
- **追蹤引擎** — 即時追蹤開信、點擊、附件下載、憑證提交、舉報
- **彈性發送** — 隨機抽樣、部門篩選、排程發送、工作時間限制、避開週末、30天冷卻期
- **發信合規** — SPF/DKIM/DMARC 檢測、按 SES/Mailgun/SMTP 政策自動限速、暫時性錯誤重試
- **報表系統** — 釣魚漏斗、部門風險排名、收件人明細、累犯追蹤、趨勢分析、CSV/PDF 匯出
- **自動化** — 定期測試排程（月/季/半年）、自動選情境、完成後自動寄送報表
- **稽核日誌** — 所有操作留痕，符合 ISO 27001 / NIST CSF 稽核要求

## 技術棧

| 層級 | 技術 |
|------|------|
| Backend | Go 1.22+, Gin, GORM v2 |
| Database | MySQL 8.0+ |
| Queue | Redis 7+ |
| Frontend | TypeScript, React 19, Vite, Ant Design 6 |
| Charts | Recharts |
| Mail | SMTP / Mailgun API / AWS SES |
| Container | Docker, Docker Compose |

## 專案結構

```
phishguard/
├── backend/
│   ├── cmd/
│   │   ├── api/            # Web API server (:8080)
│   │   ├── tracker/        # 追蹤 server (:8090)
│   │   └── worker/         # 背景 worker (mail + scheduler)
│   ├── internal/
│   │   ├── model/          # GORM models (16 models)
│   │   ├── handler/        # HTTP handlers (40+ endpoints)
│   │   ├── middleware/     # auth, tenant, audit
│   │   ├── service/        # 業務邏輯 (auth, campaign, report, plan, scheduler)
│   │   ├── repo/           # DB 存取層（強制 tenant 隔離）
│   │   ├── mailer/         # 發信抽象層 + rate limiter
│   │   └── tracker/        # 追蹤引擎
│   ├── migration/          # DB migration SQL
│   ├── config/             # 設定檔載入
│   └── Dockerfile
├── frontend/               # React SPA
│   └── Dockerfile
├── deploy/
│   ├── setup.sh            # VPS 一鍵部署腳本
│   └── nginx.conf          # Reverse proxy 設定
├── docs/
│   ├── architecture.md     # 系統架構
│   ├── phases.md           # 階段規劃與驗收條件
│   ├── frontend-design.md  # 前端 UX 設計
│   └── api-convention.md   # API 回應格式規範
├── docker-compose.yml      # 開發環境
├── docker-compose.prod.yml # 生產環境
└── README.md
```

## 快速開始（開發）

```bash
# 啟動 MySQL + Redis
docker compose up -d

# 啟動 backend
cd backend
go run ./cmd/api &
go run ./cmd/tracker &
go run ./cmd/worker &

# 啟動 frontend
cd frontend
npm install && npm run dev
```

## 生產部署（VPS）

```bash
# 1. 第一次執行 — 產生設定檔
./deploy/setup.sh

# 2. 編輯設定
vim .env.prod  # 設定 TRACKER_BASE_URL=https://t.yourdomain.com 和管理員帳號

# 3. 建置並啟動
./deploy/setup.sh

# 4. 設定 HTTPS
sudo certbot --nginx -d app.yourdomain.com -d t.yourdomain.com
```

最低需求：Ubuntu 22.04, 2GB RAM, Docker

## 文件

- [系統架構](docs/architecture.md)
- [階段規劃與驗收條件](docs/phases.md)
- [前端 UX 設計](docs/frontend-design.md)
- [API 回應格式規範](docs/api-convention.md)

## License

MIT License — see [LICENSE](LICENSE) for details.
