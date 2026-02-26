# SKILL: フロントエンド（Next.js 15 / TypeScript）

## バージョン・ランタイム

| パッケージ | バージョン |
|-----------|-----------|
| Next.js | 15.x (App Router) |
| React | 19.x |
| TypeScript | 5.x |
| Node.js | 22.x LTS |
| パッケージマネージャー | pnpm |

---

## ディレクトリ構成

```
frontend/
├── src/
│   ├── app/                    # App Router ページ
│   │   ├── layout.tsx          # ルートレイアウト
│   │   ├── page.tsx            # タイムライン（/）
│   │   ├── login/page.tsx
│   │   ├── register/page.tsx
│   │   ├── post/
│   │   │   └── [id]/page.tsx   # 俳句詳細
│   │   └── profile/
│   │       └── [username]/page.tsx
│   ├── components/
│   │   ├── ui/                 # 汎用UIコンポーネント（Button, Input等）
│   │   ├── haiku/              # 俳句関連コンポーネント
│   │   │   ├── HaikuCard.tsx   # 俳句表示カード
│   │   │   ├── HaikuForm.tsx   # 俳句投稿フォーム
│   │   │   └── HaikuCounter.tsx # 音数カウンター
│   │   ├── layout/
│   │   │   ├── Header.tsx
│   │   │   └── Footer.tsx
│   │   └── auth/
│   │       └── AuthGuard.tsx
│   ├── hooks/
│   │   ├── useHaikuValidation.ts  # 音数バリデーションフック
│   │   ├── useAuth.ts
│   │   └── useTimeline.ts
│   ├── lib/
│   │   ├── api/                # APIクライアント
│   │   │   ├── client.ts       # fetch ラッパー（認証トークン付与）
│   │   │   ├── haiku.ts
│   │   │   └── auth.ts
│   │   ├── mora.ts             # モーラカウントユーティリティ（フロント補助用）
│   │   └── constants.ts
│   ├── store/                  # Zustand ストア
│   │   ├── authStore.ts
│   │   └── timelineStore.ts
│   └── types/
│       ├── haiku.ts
│       └── user.ts
├── public/
├── .env.local                  # ローカル環境変数（gitignore）
├── .env.local.example
├── next.config.ts
├── tsconfig.json
├── tailwind.config.ts
└── package.json
```

---

## 依存パッケージ（主要）

```jsonc
{
  "dependencies": {
    "next": "^15.0.0",
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
    "zustand": "^5.0.0",          // 状態管理
    "swr": "^2.0.0",              // データフェッチング
    "tailwindcss": "^4.0.0",      // スタイリング
    "kuroshiro": "^1.2.0",        // 漢字→読み仮名変換（音数補助）
    "kuroshiro-analyzer-kuromoji": "^1.1.0"
  },
  "devDependencies": {
    "typescript": "^5.0.0",
    "@types/react": "^19.0.0",
    "eslint": "^9.0.0",
    "eslint-config-next": "^15.0.0",
    "prettier": "^3.0.0"
  }
}
```

---

## コーディング規約

### コンポーネント設計
- **Server Components をデフォルト**とし、インタラクティブな箇所のみ `"use client"` を付ける
- コンポーネントは単一責任。Props は最小限に絞る
- ファイル名: PascalCase（コンポーネント）/ camelCase（hooks, utils）

### 型定義

```typescript
// types/haiku.ts
export type HaikuPost = {
  id: string;
  ku1: string;   // 上の句（5音）
  ku2: string;   // 中の句（7音）
  ku3: string;   // 下の句（5音）
  author: User;
  likeCount: number;
  likedByMe: boolean;
  createdAt: string; // ISO 8601
};

export type CreateHaikuInput = {
  ku1: string;
  ku2: string;
  ku3: string;
};
```

### APIクライアント

```typescript
// lib/api/client.ts
// Authorization: Bearer <accessToken> を自動付与
// 401 を受け取ったらリフレッシュトークンで再試行
// 環境変数: NEXT_PUBLIC_API_URL
```

---

## 俳句投稿フォーム（HaikuForm）

### UI仕様
- 3つのテキストエリア: 上の句 / 中の句 / 下の句
- 各フィールドにリアルタイム音数カウンター表示（`HaikuCounter`）
- カウンターは補助表示のみ（正式バリデーションはバックエンド）
- 目標音数に対して 不足: グレー / 一致: グリーン / 超過: レッド で色分け

### フロント補助バリデーション（mora.ts）

```typescript
// ひらがな・カタカナ・長音符のモーラカウント
// 漢字は kuroshiro で読み仮名変換してからカウント
// カウントは補助用途のみ。バックエンドのバリデーションが正
export function countMora(text: string): number
```

---

## 状態管理（Zustand）

```typescript
// store/authStore.ts
type AuthStore = {
  user: User | null;
  accessToken: string | null;
  login: (token: string, user: User) => void;
  logout: () => void;
};

// store/timelineStore.ts
// 投稿一覧をキャッシュ。SWR と併用して無限スクロール実装
```

---

## 環境変数

```bash
# .env.local.example
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

## Lint / Format / コマンド実行

> **ローカルに Node.js / pnpm は不要。** 開発環境コンテナ内で実行する。
> 事前に `./scripts/dev-up.sh -d` で開発環境を起動しておくこと。

```bash
# ESLint
docker compose exec frontend pnpm lint

# 型チェック
docker compose exec frontend pnpm type-check

# ビルド確認（本番ビルドをコンテナ内でテスト）
docker compose exec frontend pnpm build

# package.json のスクリプト追加時
docker compose exec frontend pnpm <script-name>
```

---

## ページ一覧

| パス | 内容 | 認証 |
|------|------|------|
| `/` | タイムライン（最新俳句一覧） | 不要 |
| `/login` | ログイン | 不要 |
| `/register` | 新規登録 | 不要 |
| `/post/new` | 俳句投稿 | 必要 |
| `/post/[id]` | 俳句詳細・コメント | 不要 |
| `/profile/[username]` | プロフィール・投稿一覧 | 不要 |

---

## 注意事項

- Next.js 15 では `params` が Promise になる（`await params` が必要）
- App Router では `fetch` のキャッシュ動作に注意（`cache: 'no-store'` vs `revalidate`）
- トークンはメモリ（Zustand）に保持し、`localStorage` には保存しない（XSS対策）
- アクセストークンは短命（15分）、リフレッシュトークンは HttpOnly Cookie
