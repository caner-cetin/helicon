package lib

import (
	"github.com/caner-cetin/helicon/internal"
	"github.com/playwright-community/playwright-go"
)

// LaunchBrowser and a new Browser Context. If Chromium is not installed, it will be installed first.
// Context locale is set to `en-EN` by default.
func (h *Helicon) LaunchBrowser(
	runOptions *playwright.RunOptions,
	launchOptions *playwright.BrowserTypeLaunchOptions,
	contextOptions *playwright.BrowserNewContextOptions,
) error {
	var err error
	if runOptions == nil {
		runOptions = &playwright.RunOptions{Verbose: false, Browsers: []string{"chromium"}}
	}
	if err := playwright.Install(runOptions); err != nil {
		return err
	}
	if h.Chromium.Playwright, err = playwright.Run(); err != nil {
		return err
	}
	if launchOptions == nil {
		launchOptions = &playwright.BrowserTypeLaunchOptions{Headless: internal.Ptr(true)}
	}
	if h.Chromium.Browser, err = h.Chromium.Playwright.Chromium.Launch(*launchOptions); err != nil {
		return err
	}
	if contextOptions == nil {
		contextOptions = &playwright.BrowserNewContextOptions{BaseURL: internal.Ptr("about:blank")}
	}
	contextOptions.Locale = internal.Ptr("en-EN")
	if h.Chromium.Context, err = h.Chromium.Browser.NewContext(*contextOptions); err != nil {
		return err
	}
	return nil
}

func (h *Helicon) SetDefaultUserAgent(userAgent *string) {
	if userAgent == nil {
		h.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"
	} else {
		h.UserAgent = *userAgent
	}
}

func (h *Helicon) SpawnNewPage() (playwright.Page, error) {

	page, err := h.Chromium.Context.NewPage()
	if err != nil {
		return nil, err
	}
	if err := page.SetExtraHTTPHeaders(map[string]string{
		"User-Agent": h.UserAgent,
	}); err != nil {
		return nil, err
	}
	return page, nil
}
