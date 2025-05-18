package helicon

import (
	"encoding/base64"
	"fmt"
	"github.com/zalando/go-keyring"
	"strings"
)

// SaveTokensToKeyring under service "helicon",
//   - username => Twitter username
//   - password => all tokens concatenated with \x1F (US)
//
// order is
//   - CSRFToken (csrf) [Cookie.Raw]
//   - Auth Token [Cookie.Raw]
//   - Bearer Token	(with prefix)
func (h *Helicon) SaveTokensToKeyring() error {
	tokens := []string{
		h.Cookies.CSRFToken.Raw,
		h.Cookies.AuthToken.Raw,
		h.Cookies.BearerToken,
	}
	pass := base64.StdEncoding.EncodeToString([]byte(strings.Join(tokens, "\x1F")))
	err := keyring.Set("helicon", h.Credentials.Username, pass)
	if err != nil {
		return fmt.Errorf("failed to save tokens under service helicon with username %s: %w", h.Credentials.Username, err)
	}
	return nil
}

// LoadTokensFromKeyring into Helicon struct, accessible from
// [TwitterCredentials.CSRFToken] | [TwitterCredentials.BearerToken] | [TwitterCredentials.AuthToken] inside the struct
func (h *Helicon) LoadTokensFromKeyring() error {
	var passEncoded string
	var err error
	if passEncoded, err = keyring.Get("helicon", h.Credentials.Username); err != nil {
		return fmt.Errorf("failed to get tokens under service helicon with username %s: %w", h.Credentials.Username, err)
	}
	passDecodedBytes, err := base64.StdEncoding.DecodeString(passEncoded)
	if err != nil {
		return fmt.Errorf("failed to decode tokens: %w", err)
	}
	combinedTokens := string(passDecodedBytes)
	parts := strings.Split(combinedTokens, "\x1F")
	if len(parts) != 3 {
		return fmt.Errorf("failed to parse tokens: expected 3 parts, got %d. Raw data: '%s'", len(parts), combinedTokens)
	}
	h.Cookies.CSRFToken.Raw = parts[0]
	if err = h.Cookies.CSRFToken.Parse(); err != nil {
		return fmt.Errorf("failed to parse CSRFToken cookie: %w", err)
	}
	h.Cookies.AuthToken.Raw = parts[1]
	if err = h.Cookies.AuthToken.Parse(); err != nil {
		return fmt.Errorf("failed to parse auth_token cookie: %w", err)
	}
	h.Cookies.BearerToken = parts[2]
	return nil
}
