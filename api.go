package helicon

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func (h *Helicon) hitApi(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct GET request for url %s", url)
	}
	h.setCommonHeaders(req)
	//goland:noinspection GoLinter
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to hit %s: %w", url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("failed to close request body", "error", err)
		}
	}(resp.Body)
	var respBody []byte
	respBody, err = io.ReadAll(resp.Body)
	if resp.StatusCode > 200 {
		return nil, fmt.Errorf("unexpected status code %d from with body %s", resp.StatusCode, respBody)
	}
	return respBody, nil
}

// SetCommonHeaders including auth headers.
func (h *Helicon) setCommonHeaders(req *http.Request) {
	req.Header.Set("Authorization", h.Cookies.BearerToken)
	req.Header.Set("X-Csrf-Token", h.Cookies.CSRFToken.Value)
	req.Header.Set("X-Twitter-Auth-Type", "OAuth2Session")
	req.Header.Set("X-Twitter-Active-User", "yes")
	req.Header.Set("X-Twitter-Client-Language", "en")
	req.Header.Set("Cookie", fmt.Sprintf("auth_token=%s; ct0=%s", h.Cookies.AuthToken.Value, h.Cookies.CSRFToken.Value))
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", h.UserAgent)
}

// i dont want a separate utils file for you...
func Ptr[T any](val T) *T {
	return &val
}
