// usage: [cmd] tweet_id
// example tweet: https://x.com/gock/status/134738021831557120
// example response: https://paste.rs/IyFuJ.json
// response may vary depending on features, toggles or variables turned off
package main

import (
	"encoding/json"
	"fmt"
	"github.com/caner-cetin/helicon/lib"
	"github.com/caner-cetin/helicon/lib/models"
	"log"
	"log/slog"
	"os"
)

func main() {
	var helicon lib.Helicon
	var username = os.Getenv("HELICON_USERNAME")
	var password = os.Getenv("HELICON_PASSWORD")
	if username == "" || password == "" {
		slog.Error("username or password variable not set, please set `HELICON_USERNAME` and `HELICON_PASSWORD`")
		os.Exit(1)
	}
	helicon.SetLoginCredentials(username, password)
	helicon.SetDefaultUserAgent(nil)
	if err := helicon.LoadTokensFromKeyring(); err != nil {
		slog.Error("%s. trying to login.", err.Error())
		if err := helicon.LaunchBrowser(nil, nil, nil); err != nil {
			log.Fatal(err)
		}
		if err := helicon.Login(); err != nil {
			log.Fatal(err)
		}
	}
	var request = models.NewTweetDetailRequest(
		models.NewTweetDetailVariables(os.Args[1]),
		models.NewTweetDetailFeatures(),
		models.NewTweetDetailFieldToggles())
	resp, err := helicon.GetTweetDetails(*request)
	if err != nil {
		log.Fatal(err)
	}
	respMarshalled, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(respMarshalled))
}
