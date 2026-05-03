package collectorclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	collectordomain "repo-stat/collector/internal/domain"

	collectorrespmodel "repo-stat/collector/internal/dto"
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

func (ghac *gitHubApiClient) GetRepoInfo(ctx context.Context, owner string, repo string) (*collectordomain.RepoInfo, error) {

	url := "https://api.github.com/repos/" + owner + "/" + repo

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("request error:%w", collectordomain.ErrBadRequest)
	}

	req.Header.Set("User-Agent", "my-github-cli-tool")

	resp, err := ghac.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("github api call: %w", collectordomain.ErrInternalError)
	}

	defer func() { _ = resp.Body.Close() }()

	log.Println("rest.StatusCode", resp.StatusCode, url)
	if resp.StatusCode == http.StatusNotFound {
		return nil, collectordomain.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, collectordomain.ErrInternalError
	}

	RepoInfoSlice, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read respbody: %w", collectordomain.ErrInternalError)
	}

	var RepInfo collectorrespmodel.RepoInfo
	err = json.Unmarshal(RepoInfoSlice, &RepInfo)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse resp body: %w", collectordomain.ErrInternalError)
	}

	return &collectordomain.RepoInfo{
		FullName:    RepInfo.FullName,
		Description: RepInfo.Description,
		Forks:       RepInfo.Forks,
		Stargazers:  RepInfo.Stargazers,
		CreatedAt:   RepInfo.CreatedAt,
	}, nil
}
