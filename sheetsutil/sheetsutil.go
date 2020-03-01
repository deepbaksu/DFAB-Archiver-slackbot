package sheetsutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/slack-go/slack"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var TargetSheetId = "1jlQE7BQUYdImdn5sB8e1fohpPtMiOWKlH2euITlvZR8"

func GetOauthConfig(path string) *oauth2.Config {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v (Try getting a credentials.json from https://developers.google.com/sheets/api/quickstart/go)", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config
}

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := TokenFromFile(tokFile)
	if err != nil {
		tok = TokenFromWeb(config)
		SaveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func TokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	if err = json.NewEncoder(f).Encode(token); err != nil {
		log.Fatalf("Unable to Encode from json from path: %v, err: %v", path, err)
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
func Serialize(buf []slack.Message) [][]interface{} {
	var temp [][]interface{}

	for _, m := range buf {
		temp = append(temp, serializeMessage(m))
	}

	return temp
}
