package collectorusecase

import (
	"context"

	collectordomain "repo-stat/collector/internal/domain"
)

type collectorService struct {
	ghc collectordomain.GitHubClient
}

func NewCollectorService(ghc collectordomain.GitHubClient) *collectorService {
	return &collectorService{
		ghc: ghc,
	}
}

func (cs *collectorService) GetInfoRepo(ctx context.Context, owner string, repo string) (*collectordomain.RepoInfo, error) {

	return cs.ghc.GetRepoInfo(ctx, owner, repo)
}
