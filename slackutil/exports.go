package slackutil

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"log"
	"os"
	"strings"
)

type ChannelsValue struct {
	Channels map[string]bool
}

func (s ChannelsValue) Set(value string) error {
	words := strings.Split(value, ",")

	for _, word := range words {
		channel := strings.Trim(word, " ")
		s.Channels[channel] = true
	}

	return nil
}

func (s ChannelsValue) String() string {
	var channels []string
	for channel, ok := range s.Channels {
		if ok {
			channels = append(channels, channel)
		}
	}
	return strings.Join(channels, ",")
}

// Final export format.
type ChannelMessages struct {
	Channel string `json:"channel"`
	// unix epoch
	Datetime string         `json:"datetime"`
	Message  *slack.Message `json:"message"`
}

type ElasticSearchActionIndex struct {
	Id    string `json:"_id,omitempty"`
	Index string `json:"_index,omitempty"`
}

type ElasticSearchAction struct {
	Index  *ElasticSearchActionIndex `json:"index,omitempty"`
	Create *ElasticSearchActionIndex `json:"create,omitempty"`
	Delete *ElasticSearchActionIndex `json:"delete,omitempty"`
}

func PrintMessagesToStdoutAsNdjson(channelName string, buf []slack.Message) {
	encoder := json.NewEncoder(os.Stdout)
	indexKey := ElasticSearchAction{
		Index: &ElasticSearchActionIndex{
			Index: "slack",
		},
	}

	for _, msg := range buf {
		// Ts is the unique key
		indexKey.Index.Id = fmt.Sprintf("%v-%v", channelName, msg.Timestamp)
		if err := encoder.Encode(indexKey); err != nil {
			log.Fatalf("Failed to write indexKey(%v). See %v", indexKey, err)
		}

		channelMsg := ChannelMessages{channelName, ToUnixSeconds(&msg), &msg}
		if err := encoder.Encode(channelMsg); err != nil {
			log.Fatalf("Failed to write msg(%v). See %v", channelMsg, err)
		}
	}
}

// Extracts unix seconds from the slack message. (not safe).
func ToUnixSeconds(m *slack.Message) string {
	return strings.Split(m.Timestamp, ".")[0]
}

