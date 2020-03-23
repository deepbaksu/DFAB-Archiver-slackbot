package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/dl4ab/DFAB-Archiver-slackbot/sheetsutil"
	"github.com/dl4ab/DFAB-Archiver-slackbot/slackutil"
	"github.com/slack-go/slack"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var slackToken = flag.String("token", "", "Slack Bot Token")
var sheetId = flag.String("sheet-id", sheetsutil.TargetSheetId, "Target Sheet ID")

var beginUnixEpoch = flag.Int64("begin", time.Now().AddDate(0, 0, -1).Unix(), "The begin date for searching messages.")
var endUnixEpoch = flag.Int64("end", time.Now().Unix(), "The ending timestamp for searching messages.")

var beginTimestamp time.Time
var endTimestamp time.Time

type ChannelsValue struct {
	channels map[string]bool
}

var interestedChannels = map[string]bool{}

func (s ChannelsValue) Set(value string) error {
	words := strings.Split(value, ",")

	for _, word := range words {
		channel := strings.Trim(word, " ")
		s.channels[channel] = true
	}

	return nil
}

func (s ChannelsValue) String() string {
	var channels []string
	for channel, ok := range s.channels {
		if ok {
			channels = append(channels, channel)
		}
	}
	return strings.Join(channels, ",")
}

// Used to fetch multiple channels concurrently.
var wg sync.WaitGroup

func main() {
	defer wg.Wait()
	// Parse channels
	flag.Var(&ChannelsValue{interestedChannels}, "channels", `The names of channels to fetch (e.g., "daily_english,general")`)

	flag.Parse()

	if *slackToken == "" {
		log.Fatal("-token is not provided.")
	}

	log.Printf("Interested channels: %v", ChannelsValue{interestedChannels})

	// Try to build a Google sheets API.
	// TODO(kkweon): Refactor into sheetsutil.
	config := sheetsutil.GetOauthConfig("credentials.json")
	client := sheetsutil.GetClient(config)
	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	sheetResponse, err := srv.Spreadsheets.Get(*sheetId).Fields("sheets.properties.title").Do()
	if err != nil {
		log.Fatal(err)
	}

	existingSheetSet := map[string]bool{}
	for _, sheet := range sheetResponse.Sheets {
		title := sheet.Properties.Title
		existingSheetSet[title] = true
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
		log.Printf("Checking %v with %v", channel.Name, interestedChannels)
		ok := interestedChannels[channel.Name]
		log.Printf("ok=%v", ok)
		if ok {
			buf = slackutil.ReadMessages(api, channel.ID, historyParameters)
			log.Printf("Found messages")
			wg.Add(1)
			go writeMessagesToSheets(buf, channel, srv, existingSheetSet)
		}
	}

}

func writeMessagesToSheets(messages []slack.Message, channel slack.Channel, sheetsService *sheets.Service, existingSheetNames map[string]bool) {
	defer wg.Done()
	sheetName := channel.Name
	if _, ok := existingSheetNames[sheetName]; !ok {

		req := &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{
				&sheets.Request{
					AddSheet: &sheets.AddSheetRequest{
						Properties: &sheets.SheetProperties{
							Title: sheetName,
						},
					},
				},
			},
		}

		// Create a new sheet.
		_, err := sheetsService.Spreadsheets.BatchUpdate(*sheetId, req).Context(context.Background()).Do()

		if err != nil {
			log.Fatal(err)
		}

		// Create a new header row.
		var headers []interface{}

		for _, header := range []string{"Timestamp", "UserID", "Content"} {
			headers = append(headers, header)
		}
		_, err = sheetsService.Spreadsheets.Values.Append(*sheetId, sheetName+"!A1", &sheets.ValueRange{
			Values: [][]interface{}{headers},
		}).ValueInputOption("RAW").Do()
		if err != nil {
			log.Fatal(err)
		}
	}

	valuerrange := &sheets.ValueRange{
		Values: sheetsutil.Serialize(messages),
	}

	// Get the top most table and append to the bottom.
	_, err := sheetsService.Spreadsheets.Values.Append(*sheetId, sheetName+"!A1", valuerrange).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Error while appending values to the spreadsheet: %v", err)
	}
}
