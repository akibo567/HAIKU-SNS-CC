#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

# 本番用 secrets ファイルの存在チェック
if [[ ! -f "secrets/db_password.txt" ]]; then
  echo "❌ secrets/db_password.txt が見つかりません"
  echo "   以下を実行して作成してください:"
  echo "   mkdir -p secrets && echo 'your-strong-password' > secrets/db_password.txt"
  exit 1
fi

# 本番用 .env.prod の存在チェック
if [[ ! -f "backend/.env.prod" ]]; then
  echo "❌ backend/.env.prod が見つかりません"
  echo "   backend/.env.example を参考に作成してください"
  exit 1
fi

echo "🔨 [PROD] 本番環境をビルド中..."
docker compose -f docker-compose.prod.yml build

echo ""
echo "✅ 本番環境のビルドが完了しました"
echo "   起動するには: ./scripts/prod-up.sh"
