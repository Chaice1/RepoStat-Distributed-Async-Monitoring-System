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

func (ps *processorService) GetRepoInfo(ctx context.Context, req *processor_domain.GetRepoInfoRequestBody) (*processor_domain.RepoInfo, error) {

	fullname := req.Owner + "/" + req.Repo

	RepoInfo, err := ps.r.GetRepoInfo(ctx, fullname)

	if err != nil {
		if errors.Is(err, processor_domain.ErrNotFound) {
			err = ps.r.CreateFetchingTaskTransaction(ctx, req.Repo, req.Owner)
			if err != nil {
				return nil, err
			}
			return &processor_domain.RepoInfo{
				FullName: fullname,
				Status:   "FETCHING",
			}, nil

		}
		return nil, err
	}

	if RepoInfo.Status == "FETCHING" {
		return &processor_domain.RepoInfo{
			FullName: fullname,
			Status:   "FETCHING",
		}, nil
	}

	if RepoInfo.Status == "ERROR" {
		return nil, processor_domain.ErrNotFound

	}

	return &processor_domain.RepoInfo{
		FullName:    RepoInfo.FullName,
		Description: RepoInfo.Description,
		Stargazers:  RepoInfo.Stargazers,
		Forks:       RepoInfo.Forks,
		CreatedAt:   RepoInfo.CreatedAt,
		Status:      RepoInfo.Status,
	}, nil
}

func (ps *processorService) Ping(ctx context.Context) string {
	return "pong"
}
