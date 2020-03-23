#!/bin/sh

TIMESTAMP_FILE="./last_timestamp.txt"
BEGIN_TIMESTAMP="$(($(cat ${TIMESTAMP_FILE}) + 1))"
END_TIMESTAMP=$(date +%s)

SLACK_TOKEN="${SLACK_TOKEN:?}"
SHEET_ID="${SHEET_ID:?}"

CHANNELS="daily_english,general,test-channel-for-bots"

go run extract_messages.go \
  -token ${SLACK_TOKEN:?} \
  -begin ${BEGIN_TIMESTAMP:?} \
  -end ${END_TIMESTAMP:?} \
  -sheet-id ${SHEET_ID:?} \
  -channels ${CHANNELS:?} \
  -dry-run

if [ $? -eq 0 ]; then
  echo ${END_TIMESTAMP} > ${TIMESTAMP_FILE}
fi
