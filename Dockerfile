# PhishGuard All-in-One Container
# 包含: API Server + Tracker Server + Worker + Frontend (nginx)
# 外部需要: MySQL + Redis

# ── Stage 1: Build Go binaries ──
FROM golang:1.24-alpine AS go-builder
ENV GOTOOLCHAIN=auto
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 go build -o /bin/phishguard-api ./cmd/api && \
    CGO_ENABLED=0 go build -o /bin/phishguard-tracker ./cmd/tracker && \
    CGO_ENABLED=0 go build -o /bin/phishguard-worker ./cmd/worker

# ── Stage 2: Build Frontend ──
FROM node:20-alpine AS fe-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# ── Stage 3: Runtime ──
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata nginx supervisor

# Go binaries
COPY --from=go-builder /bin/phishguard-api /usr/local/bin/
COPY --from=go-builder /bin/phishguard-tracker /usr/local/bin/
COPY --from=go-builder /bin/phishguard-worker /usr/local/bin/

# Frontend static files
COPY --from=fe-builder /app/dist /var/www/html

# Nginx config (serves frontend + proxies /api to API server)
RUN mkdir -p /run/nginx
COPY <<'NGINX' /etc/nginx/http.d/default.conf
server {
    listen 80 default_server;

    # Frontend
    root /var/www/html;
    index index.html;
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API proxy
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        client_max_body_size 10m;
    }
}
NGINX

# Supervisord config
COPY <<'SUPERVISOR' /etc/supervisord.conf
[supervisord]
nodaemon=true
logfile=/dev/stdout
logfile_maxbytes=0

[program:api]
command=/usr/local/bin/phishguard-api
autorestart=true
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0

[program:tracker]
command=/usr/local/bin/phishguard-tracker
autorestart=true
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0

[program:worker]
command=/usr/local/bin/phishguard-worker
autorestart=true
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0

[program:nginx]
command=nginx -g "daemon off;"
autorestart=true
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0
SUPERVISOR

# Ports: 80 (frontend+API), 8090 (tracker)
EXPOSE 80 8090

CMD ["supervisord", "-c", "/etc/supervisord.conf"]
