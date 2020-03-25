# DFAB-Archiver-slackbot

[![Build Status](https://travis-ci.com/dl4ab/DFAB-Archiver-slackbot.svg?branch=kkweon-develop)](https://travis-ci.com/dl4ab/DFAB-Archiver-slackbot)
[![codecov](https://codecov.io/gh/dl4ab/DFAB-Archiver-slackbot/branch/kkweon-develop/graph/badge.svg)](https://codecov.io/gh/dl4ab/DFAB-Archiver-slackbot)
![Fetch slack messages and push.](https://github.com/dl4ab/DFAB-Archiver-slackbot/workflows/Fetch%20slack%20messages%20and%20push./badge.svg?event=schedule)

딥러닝을 공부하는 청년백수 모임 Slack의 주요 메세지를 자동으로 아카이빙하는 Slack Bot입니다.

## 프로젝트 구조

1. `[last_timestamp.txt](./last_timestamp.txt)` 는 가장 마지막으로 접한 슬랙 메시지의 시간(Unix
   Epoch)을 담고 있습니다. 그래서 두번째 잡이 실행될 때 검색할 시간 begin
   timestamp 로서 사용됩니다.
1. `[extract_messages.go](./extract_messages.go)` 는 last_timestamp.txt 에 있는
   시간부터 현재 시간까지 슬랙 메시지가 있는지 검색합니다.

```shell
❯ tree -I 'DFAB-*|*.json|*.gpg|*.mod|*.sum'
.
├── Dockerfile
├── LICENSE
├── Makefile
├── README.md
├── coverage.txt
├── elasticsearch_utils
├── extract_messages.go # main runner.
├── last_timestamp.txt # contains the last fetch run.
├── run_dev.sh
├── run_prod.sh
├── scripts # Dockerfile helper script
│   └── runner.sh
├── sheetsutil # google spreadsheet related functions
│   └── sheetsutil.go
├── slackutil  # slack related functions
│   ├── exports.go
│   ├── slacktutil_test.go
│   └── slackutil.go
└── test.log

4 directories, 15 files
```
