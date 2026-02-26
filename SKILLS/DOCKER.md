# SKILL: Docker / Nginx

## バージョン

| 技術 | バージョン |
|------|-----------|
| Docker Engine | 最新安定版 |
| Docker Compose | v2.x（`docker compose` コマンド） |
| Nginx | 1.27.x (alpine) |

---

## 環境別構成

| 環境 | ファイル | 用途 |
|------|---------|------|
| 開発 | `docker-compose.yml` | ホットリロード・デバッグ |
| 本番 | `docker-compose.prod.yml` | 最適化ビルド・セキュリティ強化 |

---

## ディレクトリ構成

```
haiku-sns-cc/
├── docker-compose.yml           # 開発環境
├── docker-compose.prod.yml      # 本番環境
├── nginx/
│   ├── nginx.dev.conf           # 開発用Nginx設定
│   └── nginx.prod.conf          # 本番用Nginx設定
├── frontend/
│   ├── Dockerfile               # 本番用（マルチステージビルド）
│   └── Dockerfile.dev           # 開発用（ホットリロード）
└── backend/
    ├── Dockerfile               # 本番用（マルチステージビルド）
    └── Dockerfile.dev           # 開発用（Air ホットリロード）
```

---

## 開発環境（docker-compose.yml）

```yaml
services:
  nginx:
    image: nginx:1.27-alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx/nginx.dev.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - frontend
      - backend

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile.dev
    volumes:
      - ./frontend:/app
      - /app/node_modules     # node_modules はコンテナ内を優先
      - /app/.next
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost/api
    expose:
      - "3000"

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.dev
    volumes:
      - ./backend:/app
    env_file:
      - ./backend/.env
    expose:
      - "8080"
    depends_on:
      db:
        condition: service_healthy

  db:
    image: postgres:17-alpine
    environment:
      POSTGRES_DB: goshichigo
      POSTGRES_USER: goshichigo
      POSTGRES_PASSWORD: goshichigo_dev
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"   # 開発時はホストから直接アクセス可
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U goshichigo -d goshichigo"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

---

## 本番環境（docker-compose.prod.yml）

```yaml
services:
  nginx:
    image: nginx:1.27-alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.prod.conf:/etc/nginx/nginx.conf:ro
      - /etc/letsencrypt:/etc/letsencrypt:ro   # TLS証明書
    restart: always
    depends_on:
      - frontend
      - backend

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
      target: runner
    environment:
      - NODE_ENV=production
      - NEXT_PUBLIC_API_URL=https://example.com/api
    expose:
      - "3000"
    restart: always

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    env_file:
      - ./backend/.env.prod
    expose:
      - "8080"
    restart: always
    depends_on:
      db:
        condition: service_healthy

  db:
    image: postgres:17-alpine
    environment:
      POSTGRES_DB: goshichigo
      POSTGRES_USER: goshichigo
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
    secrets:
      - db_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    # 本番はホストへのポート公開なし
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U goshichigo -d goshichigo"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: always

secrets:
  db_password:
    file: ./secrets/db_password.txt

volumes:
  postgres_data:
```

---

## Nginx設定

### 開発用（nginx/nginx.dev.conf）

```nginx
events {
    worker_connections 1024;
}

http {
    upstream frontend {
        server frontend:3000;
    }

    upstream backend {
        server backend:8080;
    }

    server {
        listen 80;

        # APIリクエストをバックエンドへ
        location /api/ {
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }

        # それ以外はフロントエンドへ
        location / {
            proxy_pass http://frontend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            # Next.js HMR WebSocket 対応
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
        }
    }
}
```

### 本番用（nginx/nginx.prod.conf）

```nginx
events {
    worker_connections 2048;
}

http {
    # セキュリティヘッダー
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # gzip 圧縮
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;

    upstream frontend { server frontend:3000; }
    upstream backend  { server backend:8080;  }

    # HTTP → HTTPS リダイレクト
    server {
        listen 80;
        return 301 https://$host$request_uri;
    }

    server {
        listen 443 ssl;
        ssl_certificate     /etc/letsencrypt/live/example.com/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/example.com/privkey.pem;
        ssl_protocols TLSv1.2 TLSv1.3;

        # レートリミット
        limit_req_zone $binary_remote_addr zone=api:10m rate=30r/m;

        location /api/ {
            limit_req zone=api burst=10 nodelay;
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        location / {
            proxy_pass http://frontend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }
    }
}
```

---

## Dockerfile

### バックエンド（backend/Dockerfile）

```dockerfile
# 本番用マルチステージビルド
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/server ./cmd/server

FROM scratch
COPY --from=builder /app/bin/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
```

### バックエンド開発用（backend/Dockerfile.dev）

```dockerfile
FROM golang:1.24-alpine
WORKDIR /app
# Air: ホットリロードツール
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN go mod download
EXPOSE 8080
CMD ["air"]
```

### フロントエンド（frontend/Dockerfile）

```dockerfile
# 本番用マルチステージビルド
FROM node:22-alpine AS base
RUN npm install -g pnpm

FROM base AS deps
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
ENV NEXT_TELEMETRY_DISABLED=1
RUN pnpm build

FROM base AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static
USER nextjs
EXPOSE 3000
CMD ["node", "server.js"]
```

### フロントエンド開発用（frontend/Dockerfile.dev）

```dockerfile
FROM node:22-alpine
RUN npm install -g pnpm
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN pnpm install
EXPOSE 3000
CMD ["pnpm", "dev"]
```

---

## 運用コマンド

```bash
# 開発環境起動
docker compose up --build

# 開発環境停止
docker compose down

# 本番環境起動
docker compose -f docker-compose.prod.yml up -d --build

# バックエンドのマイグレーション（開発）
docker compose exec backend migrate -path ./migrations -database $DATABASE_URL up

# ログ確認
docker compose logs -f backend
docker compose logs -f nginx

# DBに直接接続（開発）
docker compose exec db psql -U goshichigo -d goshichigo
```

---

## 注意事項

- `.env` / `.env.prod` / `secrets/` は必ず `.gitignore` に追加する
- 本番の PostgreSQL ポート（5432）は外部に公開しない
- 本番イメージは `scratch` または `distroless` を使い攻撃面を最小化する
- ボリューム `postgres_data` を削除するとDBデータが消えるため注意
- `docker compose down -v` は開発時のみ使用すること
