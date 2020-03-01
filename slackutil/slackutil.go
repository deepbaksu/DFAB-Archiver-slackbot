package slackutil

import (
	"log"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

// Returns `HistoryParameters` with Oldest and Latest set to the given time range.
func GetHistoryParams(beginTimestamp, endTimestamp time.Time) *slack.HistoryParameters {
	params := slack.NewHistoryParameters()
	params.Oldest = strconv.FormatInt(beginTimestamp.Unix(), 10)
	params.Latest = strconv.FormatInt(endTimestamp.Unix(), 10)
	return &params
}

// Returns the messages posted in channel.
func ReadMessages(api *slack.Client, channelId string, historyParams *slack.HistoryParameters,
) []slack.Message {
	history, err := api.GetChannelHistory(channelId, *historyParams)
	if err != nil {
		log.Fatal(err)
	}

	var buf []slack.Message

	for idx, message := range history.Messages {

		// Collect messages at the top level (not replies).
		if IsInterestedMessage(message) {
			buf = append(buf, message)
		}

		if history.HasMore && idx == len(history.Messages)-1 && historyParams.Latest > message.Timestamp {
			// Subtract 0.5 second to avoid fetching the same message.
			ts, err := strconv.ParseFloat(message.Timestamp, 64)
			if err != nil {
				log.Fatalf("Unable to parse float %v got an error %v", message.Timestamp, err)
			}
			historyParams.Latest =
				strconv.FormatFloat(ts-0.5, 'f', -1, 64)
			buf = append(buf, ReadMessages(api, channelId, historyParams)...)
		}
	}

	return buf
}

// Returns true if it's the top most message.
func IsInterestedMessage(m slack.Message) bool {
	return m.ParentUserId == "" && m.SubType == ""
}
