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

type RedisRepo interface {
	AddToChan(*domain.RepoInfo)
	GetRepoInfo(context.Context, *domain.GetRepoInfoReq) (*domain.RepoInfo, error)
}

type processorClient struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pc   processorpb.ProcessorClient
	rr   RedisRepo
}

func NewProcessorClient(addres string, log *slog.Logger, rr RedisRepo) (*processorClient, error) {

	conn, err := grpc.NewClient(addres, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	client := processorpb.NewProcessorClient(conn)

	return &processorClient{
		log:  log,
		conn: conn,
		pc:   client,
		rr:   rr,
	}, nil
}

func (pc *processorClient) Ping(ctx context.Context) domain.PingStatus {
	_, err := pc.pc.Ping(ctx, &processorpb.PingRequest{})
	if err != nil {
		pc.log.Error("processor ping failed", "error", err)
		return domain.PingStatusDown
	}
	return domain.PingStatusUp

}

func (pc *processorClient) GetInfoRepo(ctx context.Context, req *domain.GetRepoInfoReq) (*domain.RepoInfo, error) {

	RepoInfo, err := pc.rr.GetRepoInfo(ctx, req)

	if err == nil && RepoInfo != nil {
		return RepoInfo, nil
	}

	resp, err := pc.pc.GetInfoRepo(ctx, &processorpb.GetInfoRepoRequest{
		Repo:  req.Repo,
		Owner: req.Owner,
	})

	if err != nil {
		return nil, adapter_errors.ErrorHandleFromGRPCToDomainWithLog(err, pc.log, "GetInfoRepo")
	}

	RepoInfoResp := resp.GetRepoinfo()

	RepoInfo = &domain.RepoInfo{
		FullName:    RepoInfoResp.Fullname,
		Description: RepoInfoResp.Description,
		Forks:       RepoInfoResp.Forks,
		Status:      RepoInfoResp.Status,
		Stargazers:  RepoInfoResp.Stargazers,
		CreatedAt:   RepoInfoResp.Createdat,
	}
	if RepoInfo.Status == "READY" {
		pc.rr.AddToChan(RepoInfo)
	}

	return RepoInfo, nil
}

func (pc *processorClient) Close() error {
	return pc.conn.Close()
}
