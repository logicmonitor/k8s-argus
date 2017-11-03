#!/bin/bash

set -e

COLLECTOR_SIZE=${COLLECTOR_SIZE:-"small"}

i=$(hostname | rev | cut -d"-" -f1)
ID=$(echo "${COLLECTOR_IDS}" | cut -d',' -f $((i+1)))
HTTP_VERB='GET'
RESOURCE_PATH="/setting/collectors/${ID}/installers/Linux64"
QUERY_PARAMS="collectorSize=${COLLECTOR_SIZE}"
DATA=''
URL="https://${ACCOUNT}.logicmonitor.com/santaba/rest${RESOURCE_PATH}?${QUERY_PARAMS}"
EPOCH=$(date +%s%3N)
REQUEST_VARS="${HTTP_VERB}${EPOCH}${DATA}${RESOURCE_PATH}"
HMAC=$(echo -n "$REQUEST_VARS" | openssl dgst -sha256 -binary -hmac $ACCESS_KEY)
SIGNATURE=$(echo -n "$HMAC" | od -A n -v -t x1 | tr -d ' \n' | openssl enc -base64 -A)
AUTH="LMv1 ${ACCESS_ID}:${SIGNATURE}:${EPOCH}"
echo "${URL}"
curl \
    --retry 5 \
    --retry-delay 0 \
    -XGET \
    -H 'Content-Type: application/json' \
    -H "Authorization: ${AUTH}" \
    -o "LogicMonitorSetup.bin" \
    "${URL}"

chmod +x ./LogicMonitorSetup.bin \
    && ./LogicMonitorSetup.bin

PID=$(cat /usr/local/logicmonitor/agent/bin/logicmonitor-watchdog.wrapper.pid)
while [ ! -f /usr/local/logicmonitor/agent/logs/wrapper.log ]
do
    sleep 1
done

tail -f \
    /usr/local/logicmonitor/agent/logs/watchdog.log \
    /usr/local/logicmonitor/agent/logs/wrapper.log \
    /usr/local/logicmonitor/agent/logs/sbproxy.log &

tail --pid=${PID} -f /dev/null
