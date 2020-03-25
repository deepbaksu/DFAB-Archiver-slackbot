FROM golang:1.14 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build

##########
FROM alpine:latest
WORKDIR /app
VOLUME /app
COPY --from=builder /app/DFAB-Archiver-slackbot /app/DFAB-Archiver-slackbot

ENV SLACK_TOKEN ''
ENV BEGIN_TIMESTAMP ''
ENV END_TIMESTAMP ''
ENV SHEET_ID ''
ENV CHANNELS ''

ADD scripts/runner.sh /app/runner.sh

CMD sh ./runner.sh
