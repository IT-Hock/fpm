package utils

import (
	"context"
	"github.com/google/go-github/v49/github"
	"golang.org/x/oauth2"
)

var githubClient *github.Client = nil

func getGithubToken() (string, error) {
	token, err := GetSecret("fpm", "github_token")
	if err != nil {
		return "", err
	}

	return token, nil
}

func SetGithubToken(token string) error {
	err := SetSecret("fpm", "github_token", token, map[string]string{})
	if err != nil {
		return err
	}
	return nil
}

func GetGithubClient() *github.Client {
	if githubClient == nil {
		ok, githubApiToken := getGithubApiToken()
		if !ok {
			githubClient = github.NewClient(nil)
		} else {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: githubApiToken},
			)
			tc := oauth2.NewClient(context.Background(), ts)

			githubClient = github.NewClient(tc)
		}
	}
	return githubClient
}

func getGithubApiToken() (bool, string) {
	gitHubApiToken := GetEnv("FPM_GITHUB_TOKEN", "")
	if gitHubApiToken != "" {
		return true, gitHubApiToken
	}

	config := GetConfig()
	if config == nil {
		return false, ""
	}

	if config.GithubToken != "" {
		return true, config.GithubToken
	}

	return false, ""
}

func GithubGetUser(token string) (*github.User, error) {
	client := GetGithubClient()
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(context.Background(), ts)

		client = github.NewClient(tc)
	}

	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GithubHasRateLimit(client *github.Client) (bool, error, *github.RateLimits) {
	limits, _, err := client.RateLimits(context.Background())
	if err != nil {
		return true, err, limits
	}

	if limits.Core.Remaining == 0 || limits.Search.Remaining == 0 || limits.GraphQL.Remaining == 0 || limits.IntegrationManifest.Remaining == 0 {
		return true, &github.RateLimitError{
			Rate:     *limits.Core,
			Response: nil,
			Message:  "",
		}, limits
	}

	return false, nil, nil
}
