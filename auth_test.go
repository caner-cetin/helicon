package helicon_test

import (
	"fmt"
	"github.com/caner-cetin/helicon"
	"testing"
	"time"
)

func TestHelicon_Authenticate(t *testing.T) {
	var client helicon.Helicon
	if err := client.Authenticate(); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("auth token expired: %v \n", client.Cookies.AuthToken.Expires.Before(time.Now()))
	fmt.Printf("auth token expires at %s \n", client.Cookies.AuthToken.Expires.String())
	fmt.Printf("csrf token expired: %v \n", client.Cookies.AuthToken.Expires.Before(time.Now()))
	// or do whatever you want with it
}
