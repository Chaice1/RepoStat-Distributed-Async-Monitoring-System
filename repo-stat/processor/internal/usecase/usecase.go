package processor_usecase

import (
	"context"
	"errors"
	processor_domain "repo-stat/processor/internal/domain"
)

type processorService struct {
	r processor_domain.Repository
}

func NewProcessorService(r processor_domain.Repository) *processorService {
	return &processorService{
		r: r,
	}
}

func (ps *processorService) GetRepoInfo(ctx context.Context, req []*processor_domain.GetRepoInfoRequestBody) ([]*processor_domain.RepoInfo, error) {

	repositories := make([]*processor_domain.RepoInfo, len(req))

	for i := range req {
		fullname := req[i].Owner + "/" + req[i].Repo

		RepoInfo, err := ps.r.GetRepoInfo(ctx, fullname)

		if err != nil {
			if errors.Is(err, processor_domain.ErrNotFound) {
				err = ps.r.CreateFetchingTaskTransaction(ctx, req[i].Repo, req[i].Owner)
				if err != nil {
					return nil, err
				}
				repositories[i] = &processor_domain.RepoInfo{
					FullName: fullname,
					Status:   "FETCHING",
				}
				continue
			}
			return nil, err
		}

		if RepoInfo.Status == "FETCHING" {
			repositories[i] = &processor_domain.RepoInfo{
				FullName: fullname,
				Status:   "FETCHING",
			}
			continue
		}

		if RepoInfo.Status == "ERROR" {
			return nil, processor_domain.ErrNotFound

		}

		repositories[i] = &processor_domain.RepoInfo{
			FullName:    RepoInfo.FullName,
			Description: RepoInfo.Description,
			Stargazers:  RepoInfo.Stargazers,
			Forks:       RepoInfo.Forks,
			CreatedAt:   RepoInfo.CreatedAt,
			Status:      RepoInfo.Status,
		}
	}

	return repositories, nil
}
