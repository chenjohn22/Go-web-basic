#!/bin/sh

cd `dirname $0`
[ -f './envfile' ] && . ./envfile && ./gene-config.sh

exec "$@"

