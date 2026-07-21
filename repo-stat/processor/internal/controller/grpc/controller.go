package processor_controller

import (
	"context"
	"log/slog"
	processor_controller_errors "repo-stat/processor/internal/controller/errors"
	processor_domain "repo-stat/processor/internal/domain"
	processorpb "repo-stat/proto/processor"
)

type ProcessorService interface {
	GetRepoInfo(context.Context, *processor_domain.GetRepoInfoRequestBody) (*processor_domain.RepoInfo, error)
	Ping(context.Context) string
}

type processorController struct {
	ps ProcessorService
	processorpb.UnimplementedProcessorServer
	log *slog.Logger
}

func NewProcessorService(procserv ProcessorService, log *slog.Logger) *processorController {
	return &processorController{
		ps:  procserv,
		log: log,
	}
}

func (pc *processorController) GetInfoRepo(ctx context.Context, req *processorpb.GetInfoRepoRequest) (*processorpb.GetInfoRepoResponse, error) {
	resp, err := pc.ps.GetRepoInfo(ctx, &processor_domain.GetRepoInfoRequestBody{Repo: req.GetRepo(), Owner: req.GetOwner()})

	if err != nil {
		return nil, processor_controller_errors.HandleErrorFromDomainToGRPC(err, pc.log, "GetInfoRepo")
	}

	return &processorpb.GetInfoRepoResponse{
		Repoinfo: &processorpb.RepoInfo{
			Fullname:    resp.FullName,
			Description: resp.Description,
			Forks:       resp.Forks,
			Stargazers:  resp.Stargazers,
			Status:      resp.Status,
			Createdat:   resp.CreatedAt.String(),
		},
	}, nil
}

func (pc *processorController) Ping(ctx context.Context, req *processorpb.PingRequest) (*processorpb.PingResponse, error) {
	responce := pc.ps.Ping(ctx)
	return &processorpb.PingResponse{
		Reply: responce,
	}, nil
}
