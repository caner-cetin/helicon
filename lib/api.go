package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/caner-cetin/helicon/lib/models"
	"io"
	"log/slog"
	"net/http"
)

func (h *Helicon) GetTweetDetails(request models.TweetDetailRequest) (*models.TweetDetailResponse, error) {
	uri, err := request.GetURL()
	if err != nil {
		return nil, err
	}
	body, err := h.hitApi(*uri)
	if err != nil {
		return nil, err
	}
	var response models.TweetDetailResponse
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}
	return &response, nil
}

func (h *Helicon) hitApi(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct GET request for url %s", url)
	}
	req.Header.Set("Authorization", h.Credentials.BearerToken)
	req.Header.Set("X-Csrf-Token", h.Credentials.CSRFToken)
	req.Header.Set("X-Twitter-Auth-Type", "OAuth2Session")
	req.Header.Set("X-Twitter-Active-User", "yes")
	req.Header.Set("X-Twitter-Client-Language", "en")
	req.Header.Set("Cookie", fmt.Sprintf("auth_token=%s; ct0=%s", h.Credentials.AuthToken, h.Credentials.CSRFToken))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", h.UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to hit %s: %w", url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("failed to close request body: %w", err)
		}
	}(resp.Body)
	resp_body, err := io.ReadAll(resp.Body)
	if resp.StatusCode > 200 {
		return nil, fmt.Errorf("unexpected status code %d from with body %s", resp.StatusCode, resp_body)
	}
	return resp_body, nil
}
