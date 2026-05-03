package processor

import (
	"context"
	"log/slog"
	adapter_errors "repo-stat/api/internal/adapter/errors"
	"repo-stat/api/internal/domain"
	processorpb "repo-stat/proto/processor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type processorClient struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pc   processorpb.ProcessorClient
}

func NewProcessorClient(addres string, log *slog.Logger) (*processorClient, error) {

	conn, err := grpc.NewClient(addres, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	client := processorpb.NewProcessorClient(conn)

	return &processorClient{
		log:  log,
		conn: conn,
		pc:   client,
	}, nil
}

func (pc *processorClient) GetInfoRepo(ctx context.Context, req []*domain.GetRepoInfoReq) ([]*domain.RepoInfo, error) {

	req_to_GRPC := make([]*processorpb.GetInfoRepoRequest, len(req))

	for i := range req {
		req_to_GRPC[i] = &processorpb.GetInfoRepoRequest{
			Repo:  req[i].Repo,
			Owner: req[i].Owner,
		}
	}

	resp, err := pc.pc.GetInfoRepositories(ctx, &processorpb.GetInfoRepositoriesRequest{
		Req: req_to_GRPC,
	})

	if err != nil {
		return nil, adapter_errors.ErrorHandleFromGRPCToDomainWithLog(err, pc.log, "GetInfoRepo")
	}

	repositoriesInfo := make([]*domain.RepoInfo, len(resp.GetRepositoriesinfo()))

	for i, item := range resp.GetRepositoriesinfo() {
		repositoriesInfo[i] = &domain.RepoInfo{
			FullName:    item.Fullname,
			Description: item.Description,
			Forks:       item.Forks,
			Status:      item.Status,
			Stargazers:  item.Stargazers,
			CreatedAt:   item.Createdat,
		}
	}
	return repositoriesInfo, nil
}

func (pc *processorClient) Close() error {
	return pc.conn.Close()
}
