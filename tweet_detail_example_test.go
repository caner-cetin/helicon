package helicon_test

import (
	"encoding/json"
	"fmt"
	"github.com/caner-cetin/helicon"
	"log"
	"os"
	"testing"
)

func TestHelicon_GetTweetDetails(t *testing.T) {
	var client helicon.Helicon
	if err := client.Authenticate(); err != nil {
		log.Fatal(err)
	}
	var request = helicon.NewTweetDetailRequest(
		helicon.NewTweetDetailVariables(os.Args[1]),
		helicon.NewTweetDetailFeatures(),
		helicon.NewTweetDetailFieldToggles())
	resp, err := client.GetTweetDetails(*request)
	if err != nil {
		log.Fatal(err)
	}
	respMarshalled, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(respMarshalled))
}
