#!/usr/bin/env bash
# dump-schema.sh
# Dumps the full PostgreSQL schema (no data) to schema/schema.sql
# Usage: DATABASE_URL=postgres://user:pass@host:5432/dbname ./scripts/dump-schema.sh

set -euo pipefail

echo "setting up environment variables..."

set -a
source .env
set +a

echo "environment variables set up done"

if [[ -z "${DATABASE_URL:-}" ]]; then
  echo "❌  DATABASE_URL is not set."
  echo "    Export it first: export DATABASE_URL=postgres://user:pass@host:5432/dbname"
  exit 1
fi

mkdir -p "schema"
touch "schema/schema.sql"

echo "Dumping schema from $DATABASE_URL …"

# pg_dump flags:
#   --schema-only      → DDL only, no data
#   --no-owner         → skip OWNER TO statements (portable)
#   --no-privileges    → skip GRANT/REVOKE (portable)
#   --no-comments      → cleaner output
#   --schema=public    → only the public schema (adjust if needed)
pg_dump \
  --schema-only \
  --no-owner \
  --no-privileges \
  --no-comments \
  --schema=public \
  "$DATABASE_URL" \
  > "./schema/schema.sql"

echo "✅  Schema written to schema/schema.sql"
