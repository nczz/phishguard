#!/bin/bash
set -e

# PhishGuard VPS 一鍵部署腳本
# 需求：Ubuntu 22.04+, 2GB RAM, Docker + Docker Compose

echo "🛡️ PhishGuard 部署開始"

# ── 1. 檢查 Docker ──
if ! command -v docker &>/dev/null; then
    echo "安裝 Docker..."
    curl -fsSL https://get.docker.com | sh
    sudo usermod -aG docker $USER
    echo "Docker 已安裝，請重新登入後再執行此腳本"
    exit 1
fi

if ! command -v docker compose &>/dev/null; then
    echo "❌ 需要 Docker Compose v2"
    exit 1
fi

# ── 2. 設定環境變數 ──
if [ ! -f .env.prod ]; then
    echo "建立 .env.prod..."
    JWT_SECRET=$(openssl rand -hex 32)
    DB_PASS=$(openssl rand -hex 16)

    cat > .env.prod << EOF
# PhishGuard Production Config
DB_PASS=$DB_PASS
DB_ROOT_PASS=$(openssl rand -hex 16)
JWT_SECRET=$JWT_SECRET

# 改成你的域名
TRACKER_BASE_URL=https://t.yourdomain.com

# 初始管理員
ADMIN_EMAIL=admin@yourdomain.com
ADMIN_PASSWORD=$(openssl rand -base64 12)
EOF

    echo "✅ .env.prod 已建立，請編輯域名和管理員設定："
    echo "   vim .env.prod"
    echo ""
    cat .env.prod
    echo ""
    echo "編輯完成後重新執行此腳本"
    exit 0
fi

# ── 3. 建置並啟動 ──
echo "建置 Docker images..."
docker compose -f docker-compose.prod.yml --env-file .env.prod build

echo "啟動服務..."
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d

echo "等待服務就緒..."
sleep 10

# ── 4. 驗證 ──
echo ""
echo "=== 服務狀態 ==="
docker compose -f docker-compose.prod.yml ps

echo ""
echo "✅ PhishGuard 部署完成！"
echo ""
echo "📋 下一步："
echo "  1. 設定 DNS："
echo "     app.yourdomain.com → VPS IP"
echo "     t.yourdomain.com   → VPS IP"
echo ""
echo "  2. 設定 HTTPS（建議用 certbot）："
echo "     sudo apt install certbot python3-certbot-nginx"
echo "     sudo certbot --nginx -d app.yourdomain.com -d t.yourdomain.com"
echo ""
echo "  3. 登入管理後台："
echo "     https://app.yourdomain.com"
echo "     帳號密碼見 .env.prod"
echo ""
echo "=== 常用指令 ==="
echo "  查看日誌：docker compose -f docker-compose.prod.yml logs -f"
echo "  重啟服務：docker compose -f docker-compose.prod.yml restart"
echo "  停止服務：docker compose -f docker-compose.prod.yml down"
echo "  更新部署：git pull && ./deploy/setup.sh"
