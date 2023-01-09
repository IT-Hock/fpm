package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type DeviceAuthConfig struct {
	ClientId      string
	ClientSecret  string
	CodeEndpoint  string
	TokenEndpoint string
	Scopes        []string
	GrantType     string
}

type DeviceAuthAnswer struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationUri string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type TokenAnswer struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`

	Error string `json:"error,omitempty"`
}

type TokenAnswerError struct {
	Error string      `json:"error"`
	Token TokenAnswer `json:"token_answer"`
}

func RequestDeviceCode(client *http.Client, config *DeviceAuthConfig) (string, error) {
	// Create the request.
	req, err := http.NewRequest("POST", config.CodeEndpoint, nil)
	if err != nil {
		return "", err
	}

	// Add the query parameters.
	q := req.URL.Query()
	q.Set("client_id", config.ClientId)
	q.Set("scope", strings.Join(config.Scopes, " "))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Accept", "application/json")

	// Make the request.
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"request for device code authorisation returned status %v (%v)",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Read the response. into string
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func WaitForDeviceAuthorization(client *http.Client, config *DeviceAuthConfig, code *DeviceAuthAnswer) (string, error) {
	for {
		// Create the request.
		req, err := http.NewRequest("POST", config.TokenEndpoint, nil)

		if err != nil {
			return "", err
		}

		// Add the query parameters.
		q := req.URL.Query()
		q.Set("client_id", config.ClientId)
		q.Set("client_secret", config.ClientSecret)
		q.Set("device_code", code.DeviceCode)
		q.Set("grant_type", config.GrantType)
		req.URL.RawQuery = q.Encode()

		req.Header.Set("Accept", "application/json")

		// Make the request.
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("HTTP error %v (%v) when polling for OAuth token",
				resp.StatusCode, http.StatusText(resp.StatusCode))
		}

		byteResponse, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		// TODO: Implement errors :D
		var token TokenAnswer
		err = json.Unmarshal(byteResponse, &token)
		if err != nil {
			return "", err
		}

		switch token.Error {
		case "":
			if token.AccessToken == "" {
				panic("This should not happen")
			}
			return token.AccessToken, nil
		case "incorrect_device_code":
		case "authorization_pending":

		case "slow_down":
			code.Interval *= 2
		case "access_denied":
			return "", fmt.Errorf("access denied")
		default:
			return "", fmt.Errorf("authorization failed: %v", token.Error)
		}

		time.Sleep(time.Duration(code.Interval) * time.Second)
	}
}
