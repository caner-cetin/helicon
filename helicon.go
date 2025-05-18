package helicon

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Helicon struct {
	Credentials TwitterCredentials
	UserAgent   string
	Cookies     TwitterCookies
}
type TwitterCredentials struct {
	Username string
	Password string
}

// TwitterCookies are in full format, not value, not key=value, saved with full Set-Cookie format.
type TwitterCookies struct {
	// means CSRF token, to be used in x-csrf-token and ct0 cookie.
	CSRFToken Cookie // ct0
	AuthToken Cookie // auth_token
	// this one is a bit... strange.
	// its in `Bearer XXX` format so you dont need to append Bearer to suffix, but it is not from a cookie.
	// it is not even your own "bearer token". it is a static bearer token thats used for
	// client <-> server communication.
	// you can find this token from X's main javascript file.
	// without even logging in, you can find this value. see [Helicon.FindAnonymousBearerToken]
	//
	// no one knows when does this token expire. but at least it doesnt require login process to refresh, it is just
	// finding a token inside script file with regex.
	//
	// if you get Unauthenticated error from any of the API's, and CSRFToken // AuthToken is still valid, refresh Bearer Token with
	// [Helicon.FindAnonymousBearerToken]
	BearerToken string
}

type Cookie struct {
	// raw representation of this cookie
	// ct0=XXX; Max-Age=21600; Expires=Mon, 19 May 2025 00:42:35 GMT; Path=/; Domain=.x.com; Secure
	Raw string
	// following fields are filled with [Helicon.LoadTokensFromKeyring] using [*Cookie.Parse].
	Key         string
	Value       string
	MaxAge      int
	Expires     time.Time
	Path        string
	Domain      string
	Secure      bool
	HttpOnly    *bool
	Partitioned *bool
	SameSite    *string
}

func (c *Cookie) Parse() error {
	var err error
	for part := range strings.SplitSeq(c.Raw, ";") {
		part = strings.TrimSpace(part)
		var k, v string
		if strings.Contains(part, "=") {
			attrParts := strings.SplitN(part, "=", 2)
			k = attrParts[0]
			if len(attrParts) > 1 {
				v = attrParts[1]
			}
		} else {
			k = part
		}
		switch strings.ToLower(k) {
		case "expires":
			c.Expires, err = time.Parse(time.RFC1123, v) // common format for Expires
			if err != nil {
				// try time.RFC1123Z as another common variant
				c.Expires, err = time.Parse(time.RFC1123Z, v)
				if err != nil {
					return fmt.Errorf("cannot parse Expires %s: %w", v, err)
				}
			}
		case "path":
			c.Path = v
		case "domain":
			c.Domain = v
		case "secure":
			c.Secure = true
		case "httponly":
			c.HttpOnly = Ptr(true)
		case "partitioned":
			c.Partitioned = Ptr(true)
		case "samesite":
			c.SameSite = Ptr(v)
		case "max-age":
			c.MaxAge, err = strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("cannot parse Max-Age %s as int: %w", v, err)
			}
		default:
			c.Key, c.Value, _ = strings.Cut(part, "=")
		}
	}
	return nil
}
