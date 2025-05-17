package lib

import (
	"encoding/base64"
	"fmt"
	"github.com/zalando/go-keyring"
	"strings"
)

func (h *Helicon) SaveTokensToKeyring() error {
	tokens := []string{
		h.Credentials.CSRFToken,
		h.Credentials.BearerToken,
		h.Credentials.AuthToken,
	}
	pass := base64.StdEncoding.EncodeToString([]byte(strings.Join(tokens, "\x1F")))
	err := keyring.Set("helicon", h.Credentials.User.Username, pass)
	if err != nil {
		return fmt.Errorf("failed to save tokens under service helicon with username %s: %v", h.Credentials.User.Username, err)
	}
	return nil
}

// LoadTokensFromKeyring into Helicon struct, accessible from
// [TwitterCredentials.CSRFToken] | [TwitterCredentials.BearerToken] | [TwitterCredentials.AuthToken] inside the struct
func (h *Helicon) LoadTokensFromKeyring() error {
	var passEncoded string
	var err error
	if passEncoded, err = keyring.Get("helicon", h.Credentials.User.Username); err != nil {
		return fmt.Errorf("failed to get tokens under service helicon with username %s: %v", h.Credentials.User.Username, err)
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
	h.Credentials.CSRFToken = parts[0]
	h.Credentials.BearerToken = parts[1]
	h.Credentials.AuthToken = parts[2]
	return nil
}
