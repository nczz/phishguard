# PhishGuard

多租戶釣魚模擬測試 SaaS 平台。協助企業定期執行釣魚演練、量化員工資安風險、產出稽核報告。

## 核心功能

- **多租戶架構** — 平台管理員為企業客戶開設獨立帳號，資料完全隔離
- **情境庫** — 預建釣魚情境（信件模板 + Landing Page + 教育頁），客戶選了就能發
- **自訂模板** — 進階客戶可自行編輯 HTML 信件與 Landing Page
- **追蹤引擎** — 即時追蹤開信、點擊、附件下載、憑證提交、舉報
- **彈性發送** — 隨機抽樣、部門篩選、分散發送、定期自動測試
- **報表系統** — 釣魚漏斗、部門風險排名、累犯追蹤、趨勢分析、PDF 匯出
- **稽核日誌** — 所有操作留痕，符合 ISO 27001 / NIST CSF 稽核要求

## 技術棧

| 層級 | 技術 |
|------|------|
| Backend | Go 1.22+, Gin, GORM v2 |
| Database | MySQL 8.0+ |
| Queue | Redis 7+ (asynq) |
| Frontend | TypeScript, React 18, Vite |
| Mail | SMTP / Mailgun API / AWS SES |
| Container | Docker, Docker Compose |

## 專案結構

```
phishguard/
├── backend/
│   ├── cmd/
│   │   ├── api/            # Web API server (port 8080)
│   │   ├── tracker/        # 追蹤 server (port 8090)
│   │   └── worker/         # 背景 worker (mail + scheduler)
│   ├── internal/
│   │   ├── model/          # GORM models
│   │   ├── handler/        # HTTP handlers
│   │   ├── middleware/     # auth, tenant, audit
│   │   ├── service/        # 業務邏輯
│   │   ├── repo/           # DB 存取層（強制 tenant 隔離）
│   │   ├── mailer/         # 發信抽象層
│   │   ├── tracker/        # 追蹤引擎
│   │   └── report/         # 報表產生
│   ├── migration/          # DB migration SQL
│   └── config/             # 設定檔載入
├── frontend/               # React SPA
├── docs/                   # 設計文件
├── docker-compose.yml
└── README.md
```

## 快速開始

### 前置需求

- Go 1.22+
- Node.js 20+
- Docker & Docker Compose

### 啟動開發環境

```bash
# 啟動 MySQL + Redis
docker compose up -d mysql redis

# 啟動 backend（三個 binary）
cd backend
go run ./cmd/api
go run ./cmd/tracker
go run ./cmd/worker

# 啟動 frontend
cd frontend
npm install
npm run dev
```

### 環境變數

複製 `.env.example` 為 `.env`，填入必要設定：

```bash
cp .env.example .env
```

## 文件

- [系統架構](docs/architecture.md)
- [階段規劃與驗收條件](docs/phases.md)

## License

Proprietary — All rights reserved.
