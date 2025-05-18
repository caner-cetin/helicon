package helicon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"io"
	"log"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// FindTwitterMainJavascriptUrl finds and returns main javascript of browser client
// like this (link might be dead, filename will change whenever twitter deploys new frontend code)
//
//	https://abs.twimg.com/responsive-web/client-web-legacy/main.175fd69a.js
func (h *Helicon) FindTwitterMainJavascriptUrl() (*string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://x.com/i/flow/login/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", "https://x.com/i/flow/login/", err)
	}
	req.Header.Set("User-Agent", h.UserAgent)
	//goland:noinspection GoLinter
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", req.URL.String(), err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("failed to close request", "error", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s, status: %s", resp.Request.URL.String(), resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %w", req.URL.String(), err)
	}
	htmlContent := string(body)
	var re *regexp.Regexp
	var matches []string
	re = regexp.MustCompile(`src=["'](https://abs\.twimg\.com/responsive-web/client-web-legacy/main\.[\w.-]+\.js)["']`)
	matches = re.FindStringSubmatch(htmlContent)
	if len(matches) == 0 {
		re = regexp.MustCompile(`src=["'](https://abs\.twimg\.com/responsive-web/client-web/main\.[\w.-]+\.js)["']`)
		matches = re.FindStringSubmatch(htmlContent)
		if len(matches) == 0 {
			return nil, fmt.Errorf("failed to locate main script source inside html")
		}
	}
	var scriptUri string
	_, scriptUri, _ = strings.Cut(re.FindStringSubmatch(htmlContent)[0], "src=")
	scriptUri = strings.TrimPrefix(scriptUri, `"`)
	scriptUri = strings.TrimSuffix(scriptUri, `"`)
	return &scriptUri, nil
}

// FindAnonymousBearerToken finds and returns the preset Bearer token inside main javascript url.
// Every main JS contains a Bearer token that is used in login flows, this function finds and returns that token.
//
// See [Helicon.FindTwitterMainJavascriptUrl] for `mainJavascriptUrl` parameter.
//
// Returned value is ready to use token, with `Bearer ` prefix included, just use it like
//
//	request.Header.Set("Authorization", anonymousToken)
func (h *Helicon) FindAnonymousBearerToken() (*string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	mainScriptUrl, err := h.FindTwitterMainJavascriptUrl()
	if err != nil {
		return nil, fmt.Errorf("failed to find main script url of Twitter")
	}
	req, err := http.NewRequest("GET", *mainScriptUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", *mainScriptUrl, err)
	}
	req.Header.Set("User-Agent", h.UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", req.URL.String(), err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("failed to close request", "error", err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s, status: %s", resp.Request.URL.String(), resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %w", req.URL.String(), err)
	}
	var re *regexp.Regexp
	var matches [][]byte
	var match string
	re = regexp.MustCompile("(Bearer)(.*?)(\"|\\z)")
	matches = re.FindAll(body, -1)
	match = string(matches[len(matches)-1])
	match = strings.TrimPrefix(match, `"`)
	match = strings.TrimSuffix(match, `"`)
	return &match, nil
}

// GenerateGuestToken doesnt actually generate anything, it just requests the login page and gets the guest ID.
// Take this ID, put it into header `x-guest-token` where needed in login flow.
func (h *Helicon) GenerateGuestToken() (*string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://x.com/i/flow/login/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", "https://x.com/i/flow/login/", err)
	}
	req.Header.Set("User-Agent", h.UserAgent)
	//goland:noinspection GoLinter
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", req.URL.String(), err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("failed to close request", "error", err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s, status: %s", resp.Request.URL.String(), resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %w", req.URL.String(), err)
	}
	htmlContent := string(body)
	var re *regexp.Regexp
	var matches []string
	re = regexp.MustCompile(`document\.cookie="gt=([0-9]+)`)
	matches = re.FindStringSubmatch(htmlContent)
	return &matches[1], nil
}

type LoginFlow struct {
	FlowToken string `json:"flow_token"`
	Status    string `json:"status"`
	Subtasks  []struct {
		SubtaskId         string `json:"subtask_id"`
		JsInstrumentation struct {
			Url       string `json:"url"`
			TimeoutMs int    `json:"timeout_ms"`
			NextLink  struct {
				LinkType string `json:"link_type"`
				LinkId   string `json:"link_id"`
			} `json:"next_link"`
		} `json:"js_instrumentation"`
	} `json:"subtasks"`
	AnonymousBearerToken string
	GuestToken           string
	GuestId              string
	UserAgent            string
	Att                  string
	CFBM                 string
}

func (h *Helicon) StartLoginFlow() (*LoginFlow, error) {
	anonymousToken, err := h.FindAnonymousBearerToken()
	if err != nil {
		return nil, fmt.Errorf("failed to find anonymous bearer token: %w", err)
	}
	guestId, err := h.GenerateGuestToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate guest token: %w", err)
	}
	body := `{
	"input_flow_data": {
		"flow_context": {
			"debug_overrides": {},
			"start_location": {
				"location": "manual_link"
			}
		}
	}
}`
	req, err := http.NewRequest(http.MethodPost, "https://api.x.com/1.1/onboarding/task.json", strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to construct request %w", err)
	}
	q := req.URL.Query()
	q.Set("flow_name", "login")
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", *anonymousToken)
	req.Header.Set("x-guest-token", *guestId)
	req.Header.Set("User-Agent", h.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to hit https://api.x.com/1.1/onboarding/task.json: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("failed to close body", "error", err)
		}
	}(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to hit %s: %w", req.URL.String(), err)
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode > 200 {
		return nil, fmt.Errorf("unexpected status %d from %s, raw body: %s", resp.StatusCode, req.URL.String(), string(respBytes))
	}
	var loginFlow LoginFlow
	if err := json.NewDecoder(bytes.NewReader(respBytes)).Decode(&loginFlow); err != nil {
		return nil, fmt.Errorf("failed to decode response from %s, %w", req.URL.String(), err)
	}
	loginFlow.AnonymousBearerToken = *anonymousToken
	loginFlow.GuestToken = *guestId
	loginFlow.UserAgent = h.UserAgent
	recvCookieLines := resp.Header.Values("Set-Cookie")
	for _, cookieLine := range recvCookieLines {
		for cookie := range strings.SplitSeq(cookieLine, ";") {
			cookie = strings.TrimSpace(cookie)
			key, value, _ := strings.Cut(cookie, "=")
			switch key {
			case "att":
				loginFlow.Att = value
			case "guest_id":
				loginFlow.GuestId = value
			case "__cf_bm":
				loginFlow.CFBM = value
			}
		}
	}
	return &loginFlow, nil
}

type SubmitJSChallengeRequest struct {
	FlowToken     string                                 `json:"flow_token"`
	SubtaskInputs []SubmitJSChallengeRequestSubtaskInput `json:"subtask_inputs"`
}

type SubmitJSChallengeRequestSubtaskInput struct {
	SubtaskId         string `json:"subtask_id"`
	JsInstrumentation struct {
		Response string `json:"response"`
		Link     string `json:"link"`
	} `json:"js_instrumentation"`
}

type SubmitJSChallengeResponse struct {
	FlowToken string `json:"flow_token"`
	Status    string `json:"status"`
	// there is also a subtasks field but we dont care about it, and we dont use the subtasks in original flow,
	// outside of solving the challenge so whatever i guess

}

func (f *LoginFlow) SolveAndSubmitJSChallenge(h *Helicon) error {
	var challengeSolution *string
	var err error
	if challengeSolution, err = f.solveJSInstrumentationChallenge(h.UserAgent); err != nil {
		return fmt.Errorf("failed to solve js challenge: %w", err)
	}
	var unMarshalledChallengeSolution interface{}
	if err := json.NewDecoder(strings.NewReader(*challengeSolution)).Decode(&unMarshalledChallengeSolution); err != nil {
		return fmt.Errorf("failed to decode challenge solution, %w. raw body: %s", err, *challengeSolution)
	}
	var requestBody SubmitJSChallengeRequest
	requestBody.FlowToken = f.FlowToken
	requestBody.SubtaskInputs = []SubmitJSChallengeRequestSubtaskInput{
		{
			SubtaskId: "LoginJsInstrumentationSubtask",
			JsInstrumentation: struct {
				Response string `json:"response"`
				Link     string `json:"link"`
			}{Response: *challengeSolution, Link: "next_link"},
		},
	}
	var bodyMarshalled = bytes.NewBuffer(nil)
	if err := json.NewEncoder(bodyMarshalled).Encode(requestBody); err != nil {
		return fmt.Errorf("failed to encode request body: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.x.com/1.1/onboarding/task.json", bodyMarshalled)
	if err != nil {
		return fmt.Errorf("could not construct request: %w", err)
	}
	req.Header.Set("Authorization", f.AnonymousBearerToken)
	req.Header.Set("User-Agent", f.UserAgent)
	req.Header.Set("x-twitter-active-user", "yes")
	req.Header.Set("x-twitter-client-language", "en")
	req.Header.Set("x-guest-token", f.GuestToken)
	var cookieHeader string
	cookieHeader = fmt.Sprintf("gt=%s", f.GuestToken)
	cookieHeader = fmt.Sprintf("%s; att=%s", cookieHeader, f.Att)
	cookieHeader = fmt.Sprintf("%s; guest_id_ads=%s; guest_id_marketing=%s; guest_id=%s", cookieHeader, f.GuestId, f.GuestId, f.GuestId)
	cookieHeader = fmt.Sprintf("%s; __cf_bm=%s", cookieHeader, f.CFBM)
	req.Header.Set("Cookie", cookieHeader)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute POST: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("failed to close request body", "error", err)
		}
	}(resp.Body)
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}
	if resp.StatusCode > 200 {
		return fmt.Errorf("unexpected status code %d from %s, raw body: %s", resp.StatusCode, resp.Request.URL.String(), string(respBytes))
	}
	var challengeResponse SubmitJSChallengeResponse
	if err := json.NewDecoder(bytes.NewReader(respBytes)).Decode(&challengeResponse); err != nil {
		return fmt.Errorf("failed to decode response, %w", err)
	}
	f.FlowToken = challengeResponse.FlowToken
	return nil
}

func (f *LoginFlow) solveJSInstrumentationChallenge(userAgent string) (*string, error) {
	target := f.Subtasks[0].JsInstrumentation.Url
	resp, err := http.Get(target)
	if err != nil {
		return nil, fmt.Errorf("failed to hit %s: %w", target, err)
	}
	if resp.StatusCode > 200 {
		return nil, fmt.Errorf("unexpected status %d from %s", resp.StatusCode, target)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}(resp.Body)
	script, err := io.ReadAll(resp.Body)
	scriptContent := string(script)
	if err != nil {
		return nil, fmt.Errorf("failed to read instrumentation script: %w", err)
	}
	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.UserAgent(userAgent),
		chromedp.NoSandbox,
	)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	defer cancelAlloc()
	taskCtx, cancelTask := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancelTask()
	ctxWithTimeout, cancelTimeout := context.WithTimeout(taskCtx, 30*time.Second)
	defer cancelTimeout()
	var evaluationResult interface{}
	jsToEvaluate := fmt.Sprintf(`
		(() => {
			return new Promise((resolve, reject) => {
				let resultToCapture;
				let hasResolved = false; 


				const originalGetElementsByName = document.getElementsByName;
				document.getElementsByName = function(name) {
					if (name === 'ui_metrics') {
						return [{ 
							set value(val) {
								if (!hasResolved) {
									hasResolved = true;
									resultToCapture = val;
									resolve(resultToCapture);
								}
							},
							get value() { return undefined; }
						}];
					}
					return originalGetElementsByName.apply(this, arguments);
				};
				try {
					eval(%s);
				} catch (e) {
					if (!hasResolved) {
						hasResolved = true;
						reject("Error during eval: " + e.toString() + (e.stack ? e.stack : ''));
					}
				}
				const promiseTimeout = 5000;
				setTimeout(() => {
					if (!hasResolved) {
						hasResolved = true;
						reject("Promise timed out waiting for ui_metrics.value to be set");
					}
				}, promiseTimeout);
			});
		})();
	`, "`"+scriptContent+"`")
	actions := []chromedp.Action{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			script := `document.open(); document.write("<!DOCTYPE html><html><head></head><body></body></html>"); document.close();`
			_, exp, err := runtime.Evaluate(script).Do(ctx)
			if err != nil {
				return fmt.Errorf("failed to evaluate js challenge: %w", err)
			}
			if exp != nil {
				return fmt.Errorf("JS exception setting content: %s", exp.Exception.Description)
			}
			return nil
		}),
		chromedp.Evaluate(jsToEvaluate, &evaluationResult, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithAwaitPromise(true)
		}),
	}
	err = chromedp.Run(ctxWithTimeout, actions...)
	if err != nil {
		return nil, fmt.Errorf("chromedp: script evaluation failed: %w", err)
	}
	resultStr, ok := evaluationResult.(string)
	if !ok {
		errMsg := fmt.Sprintf("chromedp: script evaluation did not return a string, got %T value: %v", evaluationResult, evaluationResult)
		if errMap, ok := evaluationResult.(map[string]interface{}); ok {
			if desc, ok := errMap["description"].(string); ok {
				return nil, fmt.Errorf("chromedp: JS promise rejected: %s", desc)
			}
		}
		return nil, fmt.Errorf("%s", errMsg)
	}

	if resultStr == "" {
		return nil, fmt.Errorf("chromedp: script evaluation returned an empty string, expected non-empty JSON")
	}
	return &resultStr, nil
}

type SubmitUsernameRequest struct {
	FlowToken     string          `json:"flow_token"`
	SubtaskInputs []SubtaskInputs `json:"subtask_inputs"`
}
type TextData struct {
	Result string `json:"result"`
}
type ResponseData struct {
	TextData TextData `json:"text_data"`
}
type SettingResponses struct {
	Key          string       `json:"key"`
	ResponseData ResponseData `json:"response_data"`
}
type SettingsList struct {
	SettingResponses []SettingResponses `json:"setting_responses"`
	Link             string             `json:"link"`
}
type SubtaskInputs struct {
	SubtaskID    string       `json:"subtask_id"`
	SettingsList SettingsList `json:"settings_list"`
}

type SubmitUsernameResponse struct {
	FlowToken string `json:"flow_token"`
	Status    string `json:"status"`
	// we do not care about any of the subtasks, only flow token, not even status.
}

type SubmitPasswordRequest struct {
	FlowToken     string                       `json:"flow_token"`
	SubtaskInputs []SubmitPasswordSubtaskInput `json:"subtask_inputs"`
}
type EnterPassword struct {
	Password string `json:"password"`
	Link     string `json:"link"`
}
type SubmitPasswordSubtaskInput struct {
	SubtaskID     string        `json:"subtask_id"`
	EnterPassword EnterPassword `json:"enter_password"`
}

func (f *LoginFlow) SubmitUsernameAndPassword(helicon *Helicon) error {
	var submitUsernameBody SubmitUsernameRequest
	submitUsernameBody.FlowToken = f.FlowToken
	submitUsernameBody.SubtaskInputs = append(submitUsernameBody.SubtaskInputs, SubtaskInputs{
		SubtaskID: "LoginEnterUserIdentifierSSO",
		SettingsList: SettingsList{
			Link: "next_link",
			SettingResponses: []SettingResponses{
				{
					Key: "user_identifier",
					ResponseData: ResponseData{
						TextData: TextData{
							Result: helicon.Credentials.Username,
						},
					},
				},
			},
		},
	})
	var bodyMarshalled = bytes.NewBuffer(nil)
	if err := json.NewEncoder(bodyMarshalled).Encode(submitUsernameBody); err != nil {
		return fmt.Errorf("failed to marshal request submitUsernameBody: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.x.com/1.1/onboarding/task.json", bodyMarshalled)
	if err != nil {
		return fmt.Errorf("could not construct request: %w", err)
	}
	req.Header.Set("Authorization", f.AnonymousBearerToken)
	req.Header.Set("User-Agent", f.UserAgent)
	req.Header.Set("x-twitter-active-user", "yes")
	req.Header.Set("x-twitter-client-language", "en")
	req.Header.Set("x-guest-token", f.GuestToken)
	var cookieHeader string
	cookieHeader = fmt.Sprintf("gt=%s", f.GuestToken)
	cookieHeader = fmt.Sprintf("%s; att=%s", cookieHeader, f.Att)
	cookieHeader = fmt.Sprintf("%s; guest_id_ads=%s; guest_id_marketing=%s; guest_id=%s", cookieHeader, f.GuestId, f.GuestId, f.GuestId)
	req.Header.Set("Cookie", cookieHeader)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute POST: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("failed to close request body", "error", err)
		}
	}(resp.Body)
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}
	if resp.StatusCode > 200 {
		return fmt.Errorf("unexpected status code %d from %s, raw body: %s", resp.StatusCode, resp.Request.URL.String(), string(respBytes))
	}
	var submitUsernameResponse SubmitUsernameResponse
	if err := json.NewDecoder(bytes.NewReader(respBytes)).Decode(&submitUsernameResponse); err != nil {
		return fmt.Errorf("failed to decode response, %w", err)
	}
	var submitPasswordRequest SubmitPasswordRequest
	submitPasswordRequest.FlowToken = submitUsernameResponse.FlowToken
	submitPasswordRequest.SubtaskInputs = []SubmitPasswordSubtaskInput{
		{
			EnterPassword: EnterPassword{
				Link:     "next_link",
				Password: helicon.Credentials.Password,
			},
			SubtaskID: "LoginEnterPassword",
		},
	}
	bodyMarshalled = bytes.NewBuffer(nil)
	if err := json.NewEncoder(bodyMarshalled).Encode(submitPasswordRequest); err != nil {
		return fmt.Errorf("failed to marshal request submitUsernameBody: %w", err)
	}
	req, err = http.NewRequest(http.MethodPost, "https://api.x.com/1.1/onboarding/task.json", bodyMarshalled)
	if err != nil {
		return fmt.Errorf("could not construct request: %w", err)
	}
	req.Header.Set("Authorization", f.AnonymousBearerToken)
	req.Header.Set("User-Agent", f.UserAgent)
	req.Header.Set("x-twitter-active-user", "yes")
	req.Header.Set("x-twitter-client-language", "en")
	req.Header.Set("x-guest-token", f.GuestToken)
	cookieHeader = fmt.Sprintf("gt=%s", f.GuestToken)
	cookieHeader = fmt.Sprintf("%s; att=%s", cookieHeader, f.Att)
	cookieHeader = fmt.Sprintf("%s; guest_id_ads=%s; guest_id_marketing=%s; guest_id=%s", cookieHeader, f.GuestId, f.GuestId, f.GuestId)
	req.Header.Set("Cookie", cookieHeader)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute POST: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("failed to close request body", "error", err)
		}
	}(resp.Body)
	respBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}
	if resp.StatusCode > 200 {
		return fmt.Errorf("unexpected status code %d from %s, raw body: %s", resp.StatusCode, resp.Request.URL.String(), string(respBytes))
	}
	for _, rawCookie := range resp.Header.Values("Set-Cookie") {
		rawCookieSplit := strings.Split(rawCookie, ";")
		k, _, _ := strings.Cut(rawCookieSplit[0], "=")
		k = strings.TrimSpace(k)
		switch k {
		case "ct0":
			helicon.Cookies.CSRFToken.Raw = rawCookie
		case "auth_token":
			helicon.Cookies.AuthToken.Raw = rawCookie
		}
	}
	if helicon.Cookies.CSRFToken.Raw == "" {
		slog.Info("we have to execute another login flow, current attempt was successful but logged us out...")
		err = helicon.Login()
		if err != nil {
			return err
		}
	}
	return nil
}
