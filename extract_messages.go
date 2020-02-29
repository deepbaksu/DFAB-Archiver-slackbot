package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/dl4ab/DFAB-Archiver-slackbot/sheetsutil"
	"github.com/dl4ab/DFAB-Archiver-slackbot/slackutil"
	slack "github.com/slack-go/slack"
	"google.golang.org/api/sheets/v4"
)

var slackToken = flag.String("token", "", "Slack Bot Token")

var beginUnixEpoch = flag.Int64("begin", time.Now().AddDate(0, 0, -1).Unix(), "The begin date for searching messages.")
var endUnixEpoch = flag.Int64("end", time.Now().Unix(), "The ending timestamp for searching messages.")

func ioutilMustTempFile() string {
	f, err := ioutil.TempFile("", "output.*.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	return f.Name()
}

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
		if channel.Name == "daily_english" {
			buf = slackutil.ReadMessages(api, channel.ID, historyParameters)
		}
	}

	// Get credentials.json from https://developers.google.com/sheets/api/quickstart/js
	config := sheetsutil.GetOauthConfig("credentials.json")
	client := sheetsutil.GetClient(config)
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	valuerrange := &sheets.ValueRange{
		Values: serialize(buf),
	}

	// Get the top most table and append to the bottom.
	_, err = srv.Spreadsheets.Values.Append(sheetsutil.TargetSheetId, "A1", valuerrange).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Error while appending values to the spreadsheet: %v", err)
	}
}

// Serialize one message to []interface{}
func serializeMessage(m slack.Message) []interface{} {
	log.Printf("Serializing %v", m)
	ts, _ := strconv.ParseFloat(m.Timestamp, 64)
	tm := time.Unix(int64(ts), int64((ts-float64(int64(ts)))*1000))

	return []interface{}{tm.Format(time.RFC3339), m.User, m.Text}
}

// Serializes to [][]interface{} so it can be sent over the network.
func serialize(buf []slack.Message) [][]interface{} {
	var temp [][]interface{}

	for _, m := range buf {
		if isInterestedMessage(m) {
			temp = append(temp, serializeMessage(m))
		}
	}

	return temp
}

// Returns true if it's the top most message.
func isInterestedMessage(m slack.Message) bool {
	return m.ParentUserId == "" && m.SubType == ""
}
