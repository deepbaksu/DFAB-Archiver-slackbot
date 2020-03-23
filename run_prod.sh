#!/bin/sh

TIMESTAMP_FILE="$PWD/last_timestamp.txt"
BEGIN_TIMESTAMP="$(($(cat ${TIMESTAMP_FILE}) + 1))"
END_TIMESTAMP=$(date +%s)

SLACK_TOKEN="${SLACK_TOKEN:?}"
SHEET_ID="${SHEET_ID:?}"

CHANNELS="daily_english,general"

DOCKER_IMAGE=kkweon/dfab-archiver:latest
docker pull ${DOCKER_IMAGE}
docker run --env BEGIN_TIMESTAMP=${BEGIN_TIMESTAMP:?} \
--name runner \
--env END_TIMESTAMP=${END_TIMESTAMP:?} \
--env SLACK_TOKEN=${SLACK_TOKEN:?} \
--env SHEET_ID=${SHEET_ID:?} \
--env CHANNELS=${CHANNELS:?} \
-v $PWD/token.json:/app/token.json \
-v $PWD/credentials.json:/app/credentials.json \
${DOCKER_IMAGE}

if [ "$(docker inspect runner --format='{{ .State.ExitCode }}')" -eq 0 ]; then
  echo ${END_TIMESTAMP} > ${TIMESTAMP_FILE}
fi
