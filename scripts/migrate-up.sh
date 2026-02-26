#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

# オプション: ステップ数（省略時は全て UP）
STEPS="${1:-}"

echo "⬆️  [MIGRATE] マイグレーションを実行します..."

if [[ -n "$STEPS" ]]; then
  docker compose exec backend migrate -path ./migrations -database "$DATABASE_URL" up "$STEPS"
else
  docker compose exec backend migrate -path ./migrations -database "$DATABASE_URL" up
fi

echo ""
echo "✅ マイグレーション完了"
