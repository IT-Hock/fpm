package utils

import (
	"github.com/xanzy/go-gitlab"
	"strconv"
)

var gitlabClient *gitlab.Client = nil

func GetGitLabClient() (*gitlab.Client, error) {
	if gitlabClient == nil {
		ok, gitlabApiToken := getGitlabApiToken()
		if !ok {
			client, err := gitlab.NewClient("")
			if err != nil {
				return nil, err
			}
			gitlabClient = client
		} else {
			client, err := gitlab.NewClient(gitlabApiToken)
			if err != nil {
				return nil, err
			}
			gitlabClient = client
		}

	}
	return gitlabClient, nil
}

func getGitlabApiToken() (bool, string) {
	gitLabApiToken := GetEnv("FPM_GITLAB_TOKEN", "")
	if gitLabApiToken != "" {
		return true, gitLabApiToken
	}

	config := GetConfig()
	if config == nil {
		return false, ""
	}

	if config.GitlabToken != "" {
		return true, config.GitlabToken
	}

	return false, ""
}

func GitLabHasRateLimit(client *gitlab.Client) (bool, int64) {
	// Because gitlab doesn't have a rate limit endpoint, we have to use the projects endpoint
	// and check if the error is a rate limit error

	_, response, err := client.Projects.ListProjects(&gitlab.ListProjectsOptions{})
	if err != nil {
		return true, 0
	}

	rateLimitStr := response.Header.Get("RateLimit-Remaining")
	if rateLimitStr == "" {
		rateLimitStr = response.Header.Get("X-RateLimit-Remaining")
		if rateLimitStr == "" {
			rateLimitStr = "0"
		}
	}

	rateLimitResetStr := response.Header.Get("RateLimit-Reset")
	if rateLimitResetStr == "" {
		rateLimitResetStr = response.Header.Get("X-RateLimit-Reset")
		if rateLimitResetStr == "" {
			rateLimitResetStr = "0"
		}
	}

	//retryAfter := response.Header.Get("Retry-After")

	rateLimitReset, err := strconv.ParseInt(rateLimitResetStr, 10, 64)
	if err != nil {
		return true, 0
	}
	rateLimit, err := strconv.ParseInt(rateLimitStr, 10, 64)
	if err != nil {
		return true, 0
	}
	//retryTime, err := strconv.ParseInt(retryAfter, 10, 64)

	if response.StatusCode == 429 || rateLimit <= 1 {
		if err != nil {
			return true, rateLimitReset
		}
		return true, rateLimitReset
	}

	return false, 0
}
