package utils

import (
	"encoding/json"
	"fpm/src/build"
	"net/http"
)

type GithubDeviceAuthAnswer struct {
	DeviceAuthAnswer
}

type GithubTokenAnswer struct {
	TokenAnswer
}

func GithubGetDeviceCode() (*GithubDeviceAuthAnswer, error) {
	client := http.DefaultClient

	config := &DeviceAuthConfig{
		ClientId:      build.GithubClientId,
		ClientSecret:  build.GithubClientSecret,
		CodeEndpoint:  "https://github.com/login/device/code",
		TokenEndpoint: "https://github.com/login/oauth/access_token",
		Scopes:        []string{},
		GrantType:     "urn:ietf:params:oauth:grant-type:device_code",
	}

	// Get the device code.
	dcr, err := RequestGithubDeviceCode(client, config)
	if err != nil {
		return nil, err
	}

	return dcr, nil
}

func GithubGetToken(code *GithubDeviceAuthAnswer) (string, error) {
	config := &DeviceAuthConfig{
		ClientId:      build.GithubClientId,
		ClientSecret:  build.GithubClientSecret,
		CodeEndpoint:  "https://github.com/login/device/code",
		TokenEndpoint: "https://github.com/login/oauth/access_token",
		Scopes:        []string{},
		GrantType:     "urn:ietf:params:oauth:grant-type:device_code",
	}

	return WaitForGithubDeviceAuthorization(http.DefaultClient, config, code)
}

func GithubDeviceAuth(client *http.Client) (string, error) {
	config := &DeviceAuthConfig{
		ClientId:      build.GithubClientId,
		ClientSecret:  build.GithubClientSecret,
		CodeEndpoint:  "https://github.com/login/device/code",
		TokenEndpoint: "https://github.com/login/oauth/access_token",
		Scopes:        []string{},
		GrantType:     "urn:ietf:params:oauth:grant-type:device_code",
	}
	dcr, err := RequestGithubDeviceCode(client, config)
	if err != nil {
		return "", err
	}

	return WaitForGithubDeviceAuthorization(client, config, dcr)
}

func RequestGithubDeviceCode(client *http.Client, config *DeviceAuthConfig) (*GithubDeviceAuthAnswer, error) {
	resp, err := RequestDeviceCode(client, config)
	if err != nil {
		return nil, err
	}

	println(resp)

	// Unmarshal response
	var dcr GithubDeviceAuthAnswer
	err = json.Unmarshal([]byte(resp), &dcr)
	if err != nil {
		return nil, err
	}

	//

	return &dcr, nil
}

func WaitForGithubDeviceAuthorization(client *http.Client, config *DeviceAuthConfig, code *GithubDeviceAuthAnswer) (string, error) {
	answer, err := WaitForDeviceAuthorization(client, config, &code.DeviceAuthAnswer)
	if err != nil {
		return "", err
	}

	return answer, nil
}
