# Go Shichi Go! — CLAUDE.md

575の俳句しか投稿できないSNS。
ユーザーは5音・7音・5音の3句からなる俳句を投稿・閲覧・いいねできる。

---

## プロジェクト概要

| 項目 | 内容 |
|------|------|
| アプリ名 | Go Shichi Go! |
| コンセプト | 俳句（5-7-5）専用SNS |
| リポジトリルート | `/Users/akibo/mywork/haiku-sns-cc` |

---

## ディレクトリ構成

```
haiku-sns-cc/
├── CLAUDE.md
├── SKILLS/
│   ├── FRONTEND.md      # Next.js 実装規約
│   ├── BACKEND.md       # Go + PostgreSQL 実装規約
│   └── DOCKER.md        # Docker / Nginx 運用規約
├── scripts/             # 環境操作シェルスクリプト
│   ├── dev-build.sh     # 開発環境ビルド
│   ├── dev-up.sh        # 開発環境起動（-d でバックグラウンド）
│   ├── dev-down.sh      # 開発環境停止（-v でDB削除）
│   ├── prod-build.sh    # 本番環境ビルド
│   ├── prod-up.sh       # 本番環境起動（常にバックグラウンド）
│   └── prod-down.sh     # 本番環境停止
├── frontend/            # Next.js アプリ
├── backend/             # Go API サーバー
├── nginx/               # Nginx 設定
├── docker-compose.yml       # 開発環境
└── docker-compose.prod.yml  # 本番環境
```

---

## 技術スタック（最新安定版）

| レイヤー | 技術 | バージョン |
|----------|------|-----------|
| フロントエンド | Next.js (TypeScript) | 15.x |
| ランタイム | Node.js | 22.x LTS |
| バックエンド | Go | 1.24.x |
| データベース | PostgreSQL | 17.x |
| リバースプロキシ | Nginx | 1.27.x |
| コンテナ | Docker / Docker Compose | 最新安定版 |

---

## コアドメインルール

### 俳句バリデーション（最重要）
- 投稿は必ず **上の句（5音）・中の句（7音）・下の句（5音）** の3フィールド
- 音数カウントは **モーラ（拍）** 単位
  - ひらがな・カタカナ: 1文字 = 1拍（小文字「っ」「ゃ」「ゅ」「ょ」も1拍）
  - 長音符「ー」= 1拍
  - 漢字: 読み仮名に変換してからカウント（バックエンドで処理）
  - 英数字・記号: 原則使用可（1文字 = 1拍として扱う）
- バリデーションはバックエンドが正とする（フロントはUX補助）

### API設計方針
- RESTful JSON API
- `/api/v1/` プレフィックス
- 認証: JWT（アクセストークン + リフレッシュトークン）

---

## 環境操作

> **ローカルに Go / Node.js は不要。** 全コマンドは Docker コンテナ内で実行する。

### 開発環境

```bash
./scripts/dev-build.sh        # イメージをビルド
./scripts/dev-up.sh           # 起動（フォアグラウンド・ログ表示）
./scripts/dev-up.sh -d        # 起動（バックグラウンド）
./scripts/dev-down.sh         # 停止
./scripts/dev-down.sh -v      # 停止 + DBデータ削除
```

### 本番環境

```bash
./scripts/prod-build.sh       # イメージをビルド（secrets/db_password.txt 必須）
./scripts/prod-up.sh          # 起動（バックグラウンド）
./scripts/prod-down.sh        # 停止
```

### コマンド実行（開発環境が起動中であること）

```bash
# バックエンド
docker compose exec backend go test ./...
docker compose exec backend go vet ./...

# フロントエンド
docker compose exec frontend pnpm lint
docker compose exec frontend pnpm type-check

# DB マイグレーション（手動）
docker compose exec backend migrate -path ./migrations -database "$DATABASE_URL" up
```

---

## 開発ルール

- **コミット前に必ず** Docker コンテナ内でテスト・Lintを実行（上記コマンド参照）
- 環境変数は `.env` ファイルで管理。`.env.example` を必ず更新する
- マイグレーションファイルは `backend/migrations/` に連番で管理
- 詳細な実装規約は各 SKILLS ファイルを参照

---

## SKILLSファイル一覧

- [SKILLS/FRONTEND.md](SKILLS/FRONTEND.md) — Next.js コンポーネント設計・状態管理・俳句UI
- [SKILLS/BACKEND.md](SKILLS/BACKEND.md) — Go APIサーバー・DBスキーマ・俳句バリデーション
- [SKILLS/DOCKER.md](SKILLS/DOCKER.md) — Docker Compose・Nginx設定・環境別構成
