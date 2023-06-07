#!/bin/sh

export DOMAIN

envsubst '${DOMAIN}' < /conf.template > /etc/nginx/nginx.conf

exec "$@"