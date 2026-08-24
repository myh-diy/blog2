#!/bin/sh
set -eu

cd "$(dirname "$0")/.."
env_file="${ENV_FILE:-.env.production}"
compose="docker compose --env-file $env_file -f docker-compose.prod.yml"

if [ ! -f "$env_file" ]; then
    echo "Missing $env_file. Copy .env.production.example and fill it first." >&2
    exit 1
fi

set -a
. "./$env_file"
set +a
: "${DOMAIN:?DOMAIN is required in $env_file}"
: "${CERTBOT_EMAIL:?CERTBOT_EMAIL is required in $env_file}"

mkdir -p deploy/certbot/www deploy/certbot/conf

echo "Starting HTTP bootstrap for $DOMAIN..."
NGINX_TEMPLATE=http.conf.template $compose up -d blog nginx

echo "Requesting the TLS certificate..."
$compose run --rm certbot certonly \
    --webroot --webroot-path /var/www/certbot \
    --domain "$DOMAIN" \
    --email "$CERTBOT_EMAIL" \
    --agree-tos --no-eff-email --non-interactive

echo "Switching Nginx to HTTPS..."
$compose up -d --force-recreate nginx
$compose exec nginx nginx -t
$compose ps

echo "HTTPS is ready at https://$DOMAIN"
