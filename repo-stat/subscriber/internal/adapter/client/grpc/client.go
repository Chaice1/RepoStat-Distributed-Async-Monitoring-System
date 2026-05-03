package subscriber_grpc_client

import (
	"context"
	"net/http"
	subscriber_domain "repo-stat/subscriber/internal/domain"
	"time"
)

type gitHubApiClient struct {
	HttpClient *http.Client
}

func NewGitHubApiClient() *gitHubApiClient {
	return &gitHubApiClient{
		HttpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (ghac *gitHubApiClient) GetRepoInfo(ctx context.Context, owner string, repo string) error {

	url := "https://api.github.com/repos/" + owner + "/" + repo

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return subscriber_domain.ErrBadRequest
	}

	req.Header.Set("User-Agent", "my-github-cli-tool")

	resp, err := ghac.HttpClient.Do(req)

	if err != nil {
		return subscriber_domain.ErrInternalError
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return subscriber_domain.ErrNotFound
	case http.StatusInternalServerError:
		return subscriber_domain.ErrInternalError
	}

	return nil
}
