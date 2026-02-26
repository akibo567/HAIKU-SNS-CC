# SKILL: バックエンド（Go 1.24 / PostgreSQL 17）

## バージョン

| 技術 | バージョン |
|------|-----------|
| Go | 1.24.x |
| PostgreSQL | 17.x |
| migrate | golang-migrate/migrate v4 |

---

## ディレクトリ構成

```
backend/
├── cmd/
│   └── server/
│       └── main.go             # エントリーポイント
├── internal/
│   ├── handler/                # HTTPハンドラー（ルーティング）
│   │   ├── haiku.go
│   │   ├── auth.go
│   │   └── user.go
│   ├── service/                # ビジネスロジック
│   │   ├── haiku.go            # 俳句バリデーション・CRUD
│   │   ├── auth.go             # JWT発行・検証
│   │   └── user.go
│   ├── repository/             # DB操作（sqlc生成コードを使用）
│   │   ├── haiku.go
│   │   └── user.go
│   ├── middleware/
│   │   ├── auth.go             # JWT検証ミドルウェア
│   │   ├── cors.go
│   │   └── logger.go
│   ├── mora/
│   │   └── counter.go          # モーラカウント（俳句の正式バリデーション）
│   └── config/
│       └── config.go           # 環境変数読み込み
├── migrations/                 # SQLマイグレーションファイル
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   ├── 000002_create_haiku_posts.up.sql
│   └── 000002_create_haiku_posts.down.sql
├── query/                      # sqlc用SQLクエリ
│   ├── haiku.sql
│   └── user.sql
├── sqlc.yaml
├── .env                        # 環境変数（gitignore）
├── .env.example
├── go.mod
└── go.sum
```

---

## 主要依存ライブラリ

```go
// go.mod 主要依存
require (
    github.com/go-chi/chi/v5         // HTTPルーター
    github.com/golang-jwt/jwt/v5     // JWT
    github.com/jackc/pgx/v5          // PostgreSQLドライバー
    github.com/sqlc-dev/sqlc         // SQLコード生成（devtool）
    github.com/golang-migrate/migrate/v4  // DBマイグレーション
    github.com/ikawaha/kagome/v2     // 日本語形態素解析（モーラカウント）
    golang.org/x/crypto              // bcrypt（パスワードハッシュ）
    github.com/go-playground/validator/v10 // 入力バリデーション
)
```

---

## DBスキーマ

### users テーブル

```sql
-- migrations/000001_create_users.up.sql
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username    VARCHAR(30) NOT NULL UNIQUE,
    email       VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    display_name VARCHAR(50) NOT NULL,
    bio         VARCHAR(160),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### haiku_posts テーブル

```sql
-- migrations/000002_create_haiku_posts.up.sql
CREATE TABLE haiku_posts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ku1         TEXT NOT NULL,   -- 上の句（5音）
    ku2         TEXT NOT NULL,   -- 中の句（7音）
    ku3         TEXT NOT NULL,   -- 下の句（5音）
    like_count  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_haiku_posts_user_id ON haiku_posts(user_id);
CREATE INDEX idx_haiku_posts_created_at ON haiku_posts(created_at DESC);
```

### likes テーブル

```sql
CREATE TABLE likes (
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id     UUID NOT NULL REFERENCES haiku_posts(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)
);
```

---

## APIエンドポイント

### 認証

| メソッド | パス | 説明 | 認証 |
|---------|------|------|------|
| POST | `/api/v1/auth/register` | 新規登録 | 不要 |
| POST | `/api/v1/auth/login` | ログイン | 不要 |
| POST | `/api/v1/auth/refresh` | トークンリフレッシュ | Cookie |
| POST | `/api/v1/auth/logout` | ログアウト | Bearer |

### 俳句

| メソッド | パス | 説明 | 認証 |
|---------|------|------|------|
| GET | `/api/v1/posts` | タイムライン取得（カーソルページング） | 不要 |
| POST | `/api/v1/posts` | 俳句投稿 | Bearer |
| GET | `/api/v1/posts/:id` | 俳句詳細 | 不要 |
| DELETE | `/api/v1/posts/:id` | 俳句削除（本人のみ） | Bearer |
| POST | `/api/v1/posts/:id/like` | いいね | Bearer |
| DELETE | `/api/v1/posts/:id/like` | いいね取消 | Bearer |

### ユーザー

| メソッド | パス | 説明 | 認証 |
|---------|------|------|------|
| GET | `/api/v1/users/:username` | プロフィール取得 | 不要 |
| GET | `/api/v1/users/:username/posts` | 投稿一覧 | 不要 |
| PUT | `/api/v1/users/me` | プロフィール更新 | Bearer |

---

## 俳句バリデーション（最重要）

```go
// internal/mora/counter.go

// CountMora は日本語テキストのモーラ数を返す。
// 漢字は kagome で読み仮名に変換してからカウントする。
// バックエンドのカウントが唯一の正とする。
func CountMora(text string) (int, error)

// ValidateHaiku は3句のモーラ数が 5-7-5 であることを検証する。
func ValidateHaiku(ku1, ku2, ku3 string) error
```

### モーラカウントルール
| 文字種 | カウント |
|-------|---------|
| ひらがな（大文字・小文字含む） | 1文字 = 1拍 |
| カタカナ（大文字・小文字含む） | 1文字 = 1拍 |
| 長音符「ー」 | 1拍 |
| 漢字 | 読み仮名に変換後カウント |
| 英数字・記号 | 1文字 = 1拍 |
| 空白 | カウントしない |

---

## JWT設計

```go
// アクセストークン: 15分
// リフレッシュトークン: 7日間 (HttpOnly Cookie)

type Claims struct {
    UserID   string `json:"sub"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}
```

---

## レスポンス形式

```jsonc
// 成功
{
  "data": { ... },
  "meta": { "cursor": "...", "hasNext": true }  // ページング時
}

// エラー
{
  "error": {
    "code": "INVALID_MORA_COUNT",
    "message": "上の句は5音でなければなりません（現在: 6音）"
  }
}
```

### エラーコード一覧

| コード | 説明 |
|-------|------|
| `INVALID_MORA_COUNT` | 音数が5-7-5でない |
| `UNAUTHORIZED` | 認証エラー |
| `FORBIDDEN` | 権限エラー |
| `NOT_FOUND` | リソースが存在しない |
| `CONFLICT` | 重複（username/email等） |
| `VALIDATION_ERROR` | 入力バリデーションエラー |

---

## 環境変数

```bash
# .env.example
DATABASE_URL=postgres://user:password@localhost:5432/goshichigo?sslmode=disable
JWT_SECRET=your-secret-key-here
JWT_REFRESH_SECRET=your-refresh-secret-key-here
PORT=8080
ALLOWED_ORIGINS=http://localhost:3000
```

---

## 開発コマンド

> **ローカルに Go は不要。** 開発環境コンテナ内で実行する。
> 事前に `./scripts/dev-up.sh -d` で開発環境を起動しておくこと。

```bash
# テスト実行
docker compose exec backend go test ./...

# 静的解析
docker compose exec backend go vet ./...

# マイグレーション手動実行
docker compose exec backend migrate \
  -path ./migrations \
  -database "$DATABASE_URL" up

# マイグレーション ロールバック
docker compose exec backend migrate \
  -path ./migrations \
  -database "$DATABASE_URL" down 1

# ビルド確認（コンテナ内）
docker compose exec backend go build ./...

# シェルに入る（デバッグ用）
docker compose exec backend sh
```

---

## コーディング規約

- エラーは必ず呼び出し元に伝播させる（`errors.Wrap` または `fmt.Errorf("%w", err)`）
- ハンドラーはビジネスロジックを持たない。Service層に委譲する
- トランザクションが必要な操作は Repository 層でまとめる
- コンテキスト（`context.Context`）は必ず第一引数で渡す
- ログは構造化ログ（`log/slog` 標準パッケージ）を使用
