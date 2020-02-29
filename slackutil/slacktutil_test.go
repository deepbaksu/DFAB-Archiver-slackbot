package slackutil

import (
	"testing"
	"time"
)

func TestGetHistoryParams(t *testing.T) {
	begin, _ := time.Parse("2006-01-02", "2020-02-28")
	end, _ := time.Parse("2006-01-02", "2020-02-29")

	expectedOldest := "1582848000"
	expectedLatest := "1582934400"

	params := GetHistoryParams(begin, end)

	if params.Oldest != expectedOldest {
		t.Fatalf("[params.Oldest] Expected %v but received %v", expectedOldest, params.Oldest)
	}

	if params.Latest != expectedLatest {
		t.Fatalf("[params.Latest] Expected %v but received %v", expectedLatest, params.Latest)
	}
}
