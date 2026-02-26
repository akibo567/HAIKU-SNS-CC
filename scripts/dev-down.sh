#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

# -v フラグで volumes（DBデータ）も削除
VOLUMES=""
if [[ "${1:-}" == "-v" ]]; then
  VOLUMES="-v"
  echo "⚠️  DBデータ（volumes）も削除します"
fi

echo "🛑 [DEV] 開発環境を停止します..."
docker compose down $VOLUMES

echo ""
echo "✅ 停止しました"
