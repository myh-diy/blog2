#!/bin/sh
set -eu

: "${DOMAIN:?DOMAIN is required}"
template="/etc/blog-nginx/${NGINX_TEMPLATE:-https.conf.template}"

test -f "$template"
envsubst '$DOMAIN' < "$template" > /etc/nginx/conf.d/default.conf
exec nginx -g 'daemon off;'
