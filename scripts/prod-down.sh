#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

# -v フラグで volumes（DBデータ）も削除（本番では基本使わない）
VOLUMES=""
if [[ "${1:-}" == "-v" ]]; then
  VOLUMES="-v"
  echo "⚠️  [PROD] DBデータ（volumes）も削除します。本当によいですか？ (Ctrl+C でキャンセル)"
  sleep 5
fi

echo "🛑 [PROD] 本番環境を停止します..."
docker compose -f docker-compose.prod.yml down $VOLUMES

echo ""
echo "✅ 停止しました"
