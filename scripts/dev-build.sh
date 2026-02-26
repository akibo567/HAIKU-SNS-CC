#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

echo "🔨 [DEV] ビルド中..."
docker compose build

echo ""
echo "✅ 開発環境のビルドが完了しました"
echo "   起動するには: ./scripts/dev-up.sh"
