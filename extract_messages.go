package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/dl4ab/DFAB-Archiver-slackbot/sheetsutil"
	"github.com/dl4ab/DFAB-Archiver-slackbot/slackutil"
	"github.com/slack-go/slack"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var slackToken = flag.String("token", "", "Slack Bot Token")
var sheetsApiKey = flag.String("gsheet-api", "", "Google Sheets API Key")

var beginUnixEpoch = flag.Int64("begin", time.Now().AddDate(0, 0, -1).Unix(), "The begin date for searching messages.")
var endUnixEpoch = flag.Int64("end", time.Now().Unix(), "The ending timestamp for searching messages.")

var beginTimestamp time.Time
var endTimestamp time.Time

func main() {
	flag.Parse()

	if *slackToken == "" {
		log.Fatal("-token is not provided.")
	}

	beginTimestamp = time.Unix(*beginUnixEpoch, 0)
	endTimestamp = time.Unix(*endUnixEpoch, 0)

	log.Printf("beginTimestamp = %v", beginTimestamp)
	log.Printf("endTimestamp = %v", endTimestamp)

	api := slack.New(*slackToken)
	historyParameters := slackutil.GetHistoryParams(beginTimestamp, endTimestamp)
	channels, err := api.GetChannels(false)

	if err != nil {
		log.Fatal(err)
	}

	var buf []slack.Message
	for _, channel := range channels {

		// Currently, only interested in this channel.
		if channel.Name == "daily_english" {
			buf = slackutil.ReadMessages(api, channel.ID, historyParameters)
			break
		}
	}

	// Try to build a Google sheets API.
	// TODO(kkweon): Refactor into sheetsutil.
	srv, err := sheets.NewService(context.Background(), option.WithAPIKey(*sheetsApiKey))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	valuerrange := &sheets.ValueRange{
		Values: sheetsutil.Serialize(buf),
	}

	// Get the top most table and append to the bottom.
	_, err = srv.Spreadsheets.Values.Append(sheetsutil.TargetSheetId, "A1", valuerrange).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Error while appending values to the spreadsheet: %v", err)
	}
}
