package lib

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log/slog"
	"strings"
)

// Login to Twitter and save CSRF, Bearer, Auth Token to keyring. Username and password must be set
// refer to [Helicon.SetLoginCredentials]
func (h *Helicon) Login() error {
	var usernameNotSet = h.Credentials.User.Username == ""
	var passwordNotSet = h.Credentials.User.Password == ""
	if usernameNotSet && passwordNotSet {
		return fmt.Errorf("both password and username not set")
	} else if usernameNotSet {
		return fmt.Errorf("username not set")
	} else if passwordNotSet {
		return fmt.Errorf("password not set")
	}
	if len(h.Chromium.Context.Pages()) == 0 {
		slog.Warn("no pages spawned with chromium context, spawning one for login...")
		if _, err := h.SpawnNewPage(); err != nil {
			return fmt.Errorf("failed to spawn new page: %v", err)
		}
	}
	slog.Info("using first page in context for login...")
	page := h.Chromium.Context.Pages()[0]
	resp, err := page.Goto("https://x.com/i/flow/login")
	if err != nil {
		return fmt.Errorf("failed to navigate to login flow: %v", err)
	}
	if resp.Status() >= 400 {
		return fmt.Errorf("unexpected status %d from %s", resp.Status(), resp.URL())
	}
	err = page.Locator(`input[autocomplete="username"]`).Fill(h.Credentials.User.Username)
	if err != nil {
		return fmt.Errorf("failed to locate username field: %v", err)
	}
	err = page.Locator("//button[.//span[normalize-space()='Next']]").Click()
	if err != nil {
		return fmt.Errorf("failed to click next button: %v", err)
	}
	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{State: playwright.LoadStateDomcontentloaded})
	if err != nil {
		return fmt.Errorf("timeout while waiting for DOM content to load in %s, err: %v", page.URL(), err)
	}
	err = page.Locator(`input[autocomplete="current-password"]`).Fill(h.Credentials.User.Password)
	if err != nil {
		return fmt.Errorf("failed to locate password field: %v", err)
	}
	err = page.Locator(`//button[.//span[normalize-space()='Log in']]`).Click()
	if err != nil {
		return fmt.Errorf("failed to click login button: %v", err)
	}
	// let it timeout, twitter does not stop network requests, but after timeout, page will be loaded anyways.
	_ = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{State: playwright.LoadStateNetworkidle})
	cookies, err := page.Context().Cookies("https://x.com")
	if err != nil {
		return fmt.Errorf("cannot get cookies from twitter: %v", err)
	}
	for _, cookie := range cookies {
		switch cookie.Name {
		case "ct0":
			h.Credentials.CSRFToken = cookie.Value
		case "auth_token":
			h.Credentials.AuthToken = cookie.Value
		}
	}
	err = page.Route("https://x.com/**", func(r playwright.Route) {
		isApiRequest := strings.Contains(r.Request().URL(), "api")
		isFailure := r.Request().Failure() != nil
		isBearerTokenAlreadySet := h.Credentials.BearerToken != ""
		if !isApiRequest || isFailure || isBearerTokenAlreadySet {
			if err := r.Continue(); err != nil {
				slog.Warn("failed to continue request %s: %v", r.Request().URL(), err)
			}
		}
		headers, err := r.Request().AllHeaders()
		if err != nil {
			slog.Error("error getting headers for %s: %v\n", r.Request().URL(), err)
			return
		}
		authHeader, ok := headers["authorization"]
		if ok && strings.HasPrefix(authHeader, "Bearer ") {
			h.Credentials.BearerToken = authHeader
		}
	})
	if err != nil {
		return fmt.Errorf("failed to setup router for https://x.com/**: %v", err)
	}
	resp, err = page.Goto("https://x.com/home")
	if err != nil {
		return fmt.Errorf("failed to navigate to twitter home page: %v", err)
	}
	if resp.Status() >= 400 {
		return fmt.Errorf("unexpected response from %s with status %d", resp.URL(), resp.Status())
	}
	// again, just let it timeout, by the time we timeout, there will be shitton of requests to API, and bearer token
	// will be extracted anyways.
	_ = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{State: playwright.LoadStateNetworkidle})
	if err = h.SaveTokensToKeyring(); err != nil {
		return fmt.Errorf("failed to save tokens to keyring: %v", err)
	}
	return nil
}

func (h *Helicon) SetLoginCredentials(username string, password string) {
	h.Credentials.User.Username = username
	h.Credentials.User.Password = password
}
