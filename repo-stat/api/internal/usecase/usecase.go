package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type ProcessorClient interface {
	GetInfoRepo(context.Context, []*domain.GetRepoInfoReq) ([]*domain.RepoInfo, error)
}

type SubscriberClient interface {
	DeleteSubscription(context.Context, string, string) error
	CreateSubscription(context.Context, string, string) error
	GetSubscriptions(context.Context) ([]*domain.Subscription, error)
}

type ApiGatewayUsecase struct {
	pc ProcessorClient
	sc SubscriberClient
}

func NewUsecaseApiGateway(pc ProcessorClient, sc SubscriberClient) *ApiGatewayUsecase {
	return &ApiGatewayUsecase{
		pc: pc,
		sc: sc,
	}
}

func (agu *ApiGatewayUsecase) GetInfoRepo(ctx context.Context, req []*domain.GetRepoInfoReq) ([]*domain.RepoInfo, error) {
	return agu.pc.GetInfoRepo(ctx, req)
}

func (agu *ApiGatewayUsecase) DeleteSubscription(ctx context.Context, owner string, repo string) error {
	return agu.sc.DeleteSubscription(ctx, repo, owner)
}

func (agu *ApiGatewayUsecase) CreateSubscription(ctx context.Context, owner string, repo string) error {
	return agu.sc.CreateSubscription(ctx, repo, owner)
}

func (agu *ApiGatewayUsecase) GetSubscriptions(ctx context.Context) ([]*domain.Subscription, error) {
	return agu.sc.GetSubscriptions(ctx)
}
