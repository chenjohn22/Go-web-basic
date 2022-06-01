#!/bin/sh

cd `dirname $0`
CONFIG_SAMPLE_FILENAME=/opt/config/env.json.example
CONFIG_FILENAME=/opt/config/env.json

cp -rf $CONFIG_SAMPLE_FILENAME $CONFIG_FILENAME

for var in $(printenv|grep env); do
    eval "search=$(echo "$var"|tr '=' ' '|awk '{print $1}')"
    eval "replace=$(echo "$var"|tr '=' ' '|awk '{print $2}')"
    replace=${replace//\&/\\\&}
    sed -i "s|{{ ${search} }}|${replace}|g" $CONFIG_FILENAME
done