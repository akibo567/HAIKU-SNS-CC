# Go Shichi Go! 🌸

> 575の俳句しか投稿できないSNS

上の句（5音）・中の句（7音）・下の句（5音）の三句で一句。それだけ。

---

## 必要なもの

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) （Go・Node.js のローカルインストール不要）

---

## 開発環境

### 初回セットアップ

```bash
git clone <repo-url>
cd haiku-sns-cc

# イメージをビルド
./scripts/dev-build.sh

# 起動（ログをターミナルに表示）
./scripts/dev-up.sh
```

ブラウザで **http://localhost** を開く。

> 初回ビルドは依存パッケージのダウンロードがあるため数分かかります。
> バックエンド起動時にマイグレーションが自動実行されます。

### 日常の操作

```bash
# バックグラウンドで起動
./scripts/dev-up.sh -d

# 停止
./scripts/dev-down.sh

# 停止 + DB データも全削除（やり直したいとき）
./scripts/dev-down.sh -v
```

### コマンド実行

開発環境が起動中 (`dev-up.sh -d`) の状態で実行します。

```bash
# テスト（Go）
docker compose exec backend go test ./...

# 静的解析（Go）
docker compose exec backend go vet ./...

# Lint（Next.js）
docker compose exec frontend pnpm lint

# 型チェック（Next.js）
docker compose exec frontend pnpm type-check

# バックエンドのシェルに入る
docker compose exec backend sh

# フロントエンドのシェルに入る
docker compose exec frontend sh
```

### ログ確認

```bash
docker compose logs -f           # 全サービス
docker compose logs -f backend   # バックエンドのみ
docker compose logs -f frontend  # フロントエンドのみ
```

---

## 本番環境

### 初回セットアップ

**1. DB パスワードを設定する**

```bash
mkdir -p secrets
echo 'your-strong-password-here' > secrets/db_password.txt
```

**2. バックエンドの環境変数を設定する**

```bash
cp backend/.env.example backend/.env.prod
# backend/.env.prod を編集して本番用の値に変更
```

```bash
# backend/.env.prod の設定例
DATABASE_URL=postgres://goshichigo:<db_password>@db:5432/goshichigo?sslmode=disable
JWT_SECRET=<32文字以上のランダム文字列>
JWT_REFRESH_SECRET=<32文字以上の別のランダム文字列>
PORT=8080
ALLOWED_ORIGINS=https://your-domain.com
```

**3. Nginx の TLS 証明書パスを設定する**

`nginx/nginx.prod.conf` 内の `example.com` を実際のドメインに置き換えます。

**4. ビルド・起動**

```bash
./scripts/prod-build.sh
./scripts/prod-up.sh
```

### 本番の操作

```bash
# 停止
./scripts/prod-down.sh

# ログ確認
docker compose -f docker-compose.prod.yml logs -f
```

---

## プロジェクト構成

```
haiku-sns-cc/
├── scripts/                 # 環境操作スクリプト
│   ├── dev-build.sh / dev-up.sh / dev-down.sh
│   └── prod-build.sh / prod-up.sh / prod-down.sh
├── frontend/                # Next.js 15 (App Router)
├── backend/                 # Go 1.24 + PostgreSQL 17
├── nginx/                   # リバースプロキシ設定
├── docker-compose.yml       # 開発環境
└── docker-compose.prod.yml  # 本番環境
```

## 技術スタック

| レイヤー | 技術 |
|---------|------|
| フロントエンド | Next.js 15, React 19, Tailwind CSS, Zustand, SWR |
| バックエンド | Go 1.24, chi, pgx v5, golang-migrate, golang-jwt |
| データベース | PostgreSQL 17 |
| リバースプロキシ | Nginx 1.27 |

---

## API

`http://localhost/api/v1/` 以下に REST API が生えています。

| メソッド | パス | 説明 |
|---------|------|------|
| POST | `/auth/register` | 新規登録 |
| POST | `/auth/login` | ログイン |
| GET | `/posts` | タイムライン取得 |
| POST | `/posts` | 俳句投稿 ★要認証 |
| POST | `/posts/:id/like` | いいね ★要認証 |
| GET | `/users/:username` | プロフィール取得 |

詳細は [SKILLS/BACKEND.md](SKILLS/BACKEND.md) を参照。
