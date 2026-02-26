#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

# オプション: ステップ数（省略時は 1 ステップ DOWN）
STEPS="${1:-1}"

echo "⬇️  [MIGRATE] マイグレーションをロールバックします（${STEPS}ステップ）..."

docker compose exec backend migrate -path ./migrations -database "$DATABASE_URL" down "$STEPS"

echo ""
echo "✅ ロールバック完了"
