# PhishGuard All-in-One Container
# 包含: API Server + Tracker Server + Worker + Frontend (nginx)
# 外部需要: MySQL + Redis

# ── Stage 1: Build Go binaries ──
FROM golang:1.25-alpine AS go-builder
ENV GOTOOLCHAIN=auto
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/phishguard-api ./cmd/api && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/phishguard-tracker ./cmd/tracker && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/phishguard-worker ./cmd/worker

# ── Stage 2: Build Frontend ──
FROM node:20-alpine AS fe-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# ── Stage 3: Runtime ──
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata nginx supervisor mariadb-client
ENV TZ=UTC

# Go binaries
COPY --from=go-builder /bin/phishguard-api /usr/local/bin/
COPY --from=go-builder /bin/phishguard-tracker /usr/local/bin/
COPY --from=go-builder /bin/phishguard-worker /usr/local/bin/

# Migration SQL
COPY backend/migration/ /migration/

# Frontend static files
COPY --from=fe-builder /app/dist /var/www/html

# Nginx config (serves frontend + proxies /api to API server)
RUN mkdir -p /run/nginx
COPY deploy/nginx-aio.conf /etc/nginx/http.d/default.conf

# Supervisord config
COPY deploy/supervisord.conf /etc/supervisord.conf

# Startup script: run migration then start services
COPY deploy/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Ports: 80 (frontend+API), 8090 (tracker)
EXPOSE 80 8090

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://127.0.0.1/ > /dev/null || exit 1

LABEL org.opencontainers.image.source="https://github.com/nczz/phishguard" \
      org.opencontainers.image.description="PhishGuard - 企業釣魚模擬測試平台" \
      org.opencontainers.image.licenses="MIT"

CMD ["/entrypoint.sh"]
