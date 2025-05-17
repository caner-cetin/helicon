package lib

import "github.com/playwright-community/playwright-go"

type Helicon struct {
	Credentials TwitterCredentials
	Chromium    ChromiumContext
	UserAgent   string
}

type ChromiumContext struct {
	Playwright *playwright.Playwright
	Browser    playwright.Browser
	Context    playwright.BrowserContext
	PageConfig ChromiumPageConfig
}

type ChromiumPageConfig struct {
}

type TwitterCredentials struct {
	User        TwitterUserCredentials
	CSRFToken   string
	BearerToken string
	AuthToken   string
}

type TwitterUserCredentials struct {
	Username string
	Password string
}
