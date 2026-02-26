#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

# 本番用 secrets ファイルの存在チェック
if [[ ! -f "secrets/db_password.txt" ]]; then
  echo "❌ secrets/db_password.txt が見つかりません"
  echo "   mkdir -p secrets && echo 'your-strong-password' > secrets/db_password.txt"
  exit 1
fi

if [[ ! -f "backend/.env.prod" ]]; then
  echo "❌ backend/.env.prod が見つかりません"
  exit 1
fi

echo "🚀 [PROD] 本番環境を起動します..."
docker compose -f docker-compose.prod.yml up -d

echo ""
echo "✅ 本番環境が起動しました（バックグラウンド）"
echo "   URL: http://localhost  (TLS設定済みの場合は https://)"
echo "   ログ確認: docker compose -f docker-compose.prod.yml logs -f"
echo "   停止するには: ./scripts/prod-down.sh"
