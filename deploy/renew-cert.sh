#!/bin/sh
set -eu

cd "$(dirname "$0")/.."
env_file="${ENV_FILE:-.env.production}"
compose="docker compose --env-file $env_file -f docker-compose.prod.yml"

$compose run --rm certbot renew --webroot --webroot-path /var/www/certbot --quiet
$compose exec nginx nginx -s reload
