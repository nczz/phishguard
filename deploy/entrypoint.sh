#!/bin/sh
cat > /var/www/html/config.js <<EOF
window.__CONFIG__ = {
  APP_DESCRIPTION: "${APP_DESCRIPTION:-企業釣魚模擬測試平台}"
};
EOF

# Auto-migrate: run SQL if tables don't exist
if [ -f /migration/001_initial_schema.sql ]; then
  echo "Checking database schema..."
  RESULT=$(mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASS}" "${DB_NAME}" -e "SELECT 1 FROM tenants LIMIT 1" 2>/dev/null)
  if [ -z "$RESULT" ]; then
    echo "Running database migration..."
    mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASS}" "${DB_NAME}" < /migration/001_initial_schema.sql
    echo "Migration complete."
  else
    echo "Schema already exists, skipping migration."
  fi
fi

exec supervisord -c /etc/supervisord.conf
