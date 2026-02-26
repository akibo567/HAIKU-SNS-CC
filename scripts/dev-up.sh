#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

# -d フラグでバックグラウンド起動
DETACH=""
if [[ "${1:-}" == "-d" ]]; then
  DETACH="-d"
fi

echo "🚀 [DEV] 開発環境を起動します..."
docker compose up $DETACH

if [[ -n "$DETACH" ]]; then
  echo ""
  echo "✅ バックグラウンドで起動しました"
  echo "   URL: http://localhost"
  echo "   ログ確認: docker compose logs -f"
  echo "   停止するには: ./scripts/dev-down.sh"
fi
