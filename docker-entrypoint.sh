#!/bin/sh
set -e

# Default DB path 
DB="${CFG_DB_FILENAME:-/data/app.db}"

# Remove existing DB file if present to recreate from scratch each run
if [ -f "$DB" ]; then
  echo "[entrypoint] Removing existing database file: $DB"
  rm -f "$DB"
else
  echo "[entrypoint] No existing database file at: $DB"
fi

# Execute the passed command (should be the webapi binary)
exec "$@"
