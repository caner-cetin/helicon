package helicon

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

func (h *Helicon) Login() error {
	var flow *LoginFlow
	var err error
	if flow, err = h.StartLoginFlow(); err != nil {
		return fmt.Errorf("failed to start login flow: %w", err)
	}
	if err = flow.SolveAndSubmitJSChallenge(h); err != nil {
		return fmt.Errorf("failed to submit JS challenge: %w", err)
	}
	if err = flow.SubmitUsernameAndPassword(h); err != nil {
		return fmt.Errorf("failed to submit username: %w", err)
	}
	if err = h.SaveTokensToKeyring(); err != nil {
		return fmt.Errorf("failed to save the token to keyring: %w", err)
	}
	_ = h.LoadTokensFromKeyring()

	return nil
}

func (h *Helicon) SetLoginCredentials(username string, password string) {
	h.Credentials.Username = username
	h.Credentials.Password = password
}
func (h *Helicon) Authenticate() error {
	var username = os.Getenv("HELICON_USERNAME")
	if username == "" {
		return fmt.Errorf("HELICON_USERNAME environment variable not set, cannot proceed")
	}
	var password = os.Getenv("HELICON_PASSWORD")
	if password == "" {
		return fmt.Errorf("HELICON_PASSWORD environment variable not set, cannot proceed")
	}
	h.SetLoginCredentials(username, password)
	h.SetDefaultUserAgent(nil)
	var forceLogin bool
	var err error
	forceLoginString := os.Getenv("HELICON_FORCE_LOGIN")
	if forceLoginString != "" {
		forceLogin, err = strconv.ParseBool(forceLoginString)
		if err != nil {
			slog.Warn("invalid HELICON_FORCE_LOGIN, excepted bool", "received", forceLoginString)
			slog.Warn("defaulting back to HELICON_FORCE_LOGIN=false")
			forceLogin = false
		}
		if forceLogin {
			if err = h.Login(); err != nil {
				return err
			}
		}
	}
	if err := h.LoadTokensFromKeyring(); err != nil {
		if err = h.Login(); err != nil {
			return err
		}
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
