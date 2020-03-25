package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/dl4ab/DFAB-Archiver-slackbot/sheetsutil"
	"github.com/dl4ab/DFAB-Archiver-slackbot/slackutil"
	"github.com/slack-go/slack"
	"google.golang.org/api/sheets/v4"
)

var slackToken = flag.String("token", "", "Slack Bot Token")
var sheetId = flag.String("sheet-id", sheetsutil.TargetSheetId, "Target Sheet ID")

var beginUnixEpoch = flag.Int64("begin", time.Now().AddDate(0, 0, -1).Unix(), "The begin date for searching messages.")
var endUnixEpoch = flag.Int64("end", time.Now().Unix(), "The ending timestamp for searching messages.")

var dryRun = flag.Bool("dry-run", false, "Dry run if true.")

var beginTimestamp time.Time
var endTimestamp time.Time

var interestedChannels = map[string]bool{}

// Used to fetch multiple channels concurrently.
var wg sync.WaitGroup

func main() {
	defer wg.Wait()
	// Parse channels
	flag.Var(&slackutil.ChannelsValue{Channels: interestedChannels},
		"channels", `The names of channels to fetch (e.g., "daily_english,general")`)

	flag.Parse()

	if *slackToken == "" {
		log.Fatal("-token is not provided.")
	}

	log.Printf("Interested channels: %v",
		slackutil.ChannelsValue{Channels: interestedChannels})

	// Try to build a Google sheets API.
	srv, err := sheetsutil.GetService("credentials.json")

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Get existing sheet names to decide whether to create a new sheet later.
	sheetNamesSet := sheetsutil.GetSheetNamesSet(srv, *sheetId)

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

	if len(channels) == 0 {
		log.Printf("The number of channels returned from the server is 0. Please check the SLACK API key. Subscribed channels are %v",
			slackutil.ChannelsValue{Channels: interestedChannels})
	}

	var buf []slack.Message
	for _, channel := range channels {

		// Currently, only interested in this channel.
		log.Printf("Checking %v with %v", channel.Name, interestedChannels)
		ok := interestedChannels[channel.Name]
		log.Printf("ok=%v", ok)
		if ok {
			buf = slackutil.ReadMessages(api, channel.ID, historyParameters)

			// These messages will be stored in temp.log and then it will be sent to elasticsearch.
			slackutil.PrintMessagesToStdoutAsNdjson(channel.Name, buf)

			if *dryRun {
				log.Println(buf)
				continue
			}

			if len(buf) == 0 {
				log.Printf("There is no message returned in this channel(%v).", channel.Name)
			}

			// This will send messages to the spreadsheets.
			wg.Add(1)
			go writeMessagesToSheets(buf, channel, *sheetId, srv, sheetNamesSet)

		}
	}
}

func writeMessagesToSheets(messages []slack.Message, channel slack.Channel, sheetId string, sheetsService *sheets.Service, existingSheetNames map[string]bool) {
	defer wg.Done()
	sheetName := channel.Name
	if _, ok := existingSheetNames[sheetName]; !ok {
		sheetsutil.CreateNewSheet(sheetId, sheetName, sheetsService)
	}

	data := &sheets.ValueRange{
		Values: sheetsutil.Serialize(messages),
	}

	// Get the top most table and append to the bottom.
	_, err := sheetsService.Spreadsheets.Values.Append(sheetId, sheetName+"!A1", data).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Failed to append data into the sheet. See error %v", err)
	}
}
