#!/bin/sh
./DFAB-Archiver-slackbot -token ${SLACK_TOKEN:?} -begin ${BEGIN_TIMESTAMP} -end ${END_TIMESTAMP:?} -sheet-id ${SHEET_ID:?} -channels ${CHANNELS:?} 1> temp.log
