package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/nlopes/slack"
)

func main() {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	today := time.Now()
	yesterday := today.AddDate(-1, 0, 0)
	todayTimestamp := today.UnixNano() / 1000000000
	yesterdayTimestamp := yesterday.UnixNano() / 1000000000
	todayString := strconv.FormatInt(todayTimestamp, 10)
	yesterdayString := strconv.FormatInt(yesterdayTimestamp, 10)

	historyParameters := slack.NewHistoryParameters()
	historyParameters.Latest = todayString
	historyParameters.Oldest = yesterdayString
	historyParameters.Count = 1000

	fmt.Println("date & convert Check")
	fmt.Println(today)
	fmt.Println(yesterday)

	fmt.Println(todayString)
	fmt.Println(yesterdayString)

	fmt.Println("\nHistoryParameters check")
	fmt.Println(historyParameters)

	fmt.Println("\nType Check")
	fmt.Println(reflect.TypeOf(todayString))
	fmt.Println(reflect.TypeOf(yesterdayString))
	fmt.Println(reflect.TypeOf(historyParameters.Count))

	channels, err := api.GetChannels(false)

	if err != nil {
		fmt.Printf("error : %s\n", err)
		return
	}

	for _, channel := range channels {
		if channel.Name == "daily-logs-bravo" {
			fmt.Println(channel.ID)
			history, err := api.GetChannelHistory(channel.ID, historyParameters)

			if err != nil {
				fmt.Printf("error : %s\n", err)
				return
			}

			fmt.Println("HISTORY LATEST")
			fmt.Println(history.Latest)
			fmt.Println("--------------------------------------------------------------------------------------")
			fmt.Println("HISTORY HASMORE")
			fmt.Println(history.HasMore)
			fmt.Println("--------------------------------------------------------------------------------------")
			fmt.Println("HISTORY MESSAGE")
			for idx, message := range history.Messages {
				fmt.Println(idx, message)
				fmt.Println("--------------------------------------------------------------------------------------")
			}
		}
	}
}
