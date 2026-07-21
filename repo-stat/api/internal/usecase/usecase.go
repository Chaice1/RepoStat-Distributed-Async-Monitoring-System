package usecase

import (
	"context"
	"repo-stat/api/internal/domain"

	"golang.org/x/sync/errgroup"
)

type ProcessorClient interface {
	Ping(context.Context) domain.PingStatus
	GetInfoRepo(context.Context, *domain.GetRepoInfoReq) (*domain.RepoInfo, error)
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

func (agu *ApiGatewayUsecase) Ping(ctx context.Context) domain.PingStatus {
	return agu.pc.Ping(ctx)
}

func (agu *ApiGatewayUsecase) GetInfoRepo(ctx context.Context, repo string, owner string) (*domain.RepoInfo, error) {
	return agu.pc.GetInfoRepo(ctx, &domain.GetRepoInfoReq{
		Repo:  repo,
		Owner: owner,
	})
}
func (agu *ApiGatewayUsecase) GetInfoRepositories(ctx context.Context) ([]*domain.RepoInfo, error) {

	subscriptions, err := agu.sc.GetSubscriptions(ctx)

	if err != nil {
		return nil, err
	}

	if len(subscriptions) == 0 {
		return []*domain.RepoInfo{}, nil
	}

	RepositoriesInfo := make([]*domain.RepoInfo, len(subscriptions))

	errGroup, cttx := errgroup.WithContext(ctx)

	errGroup.SetLimit(10)
	for i := range subscriptions {

		idx := i
		sub := subscriptions[idx]
		errGroup.Go(func() error {

			RepoInfo, err := agu.pc.GetInfoRepo(cttx, &domain.GetRepoInfoReq{
				Repo:  sub.RepoName,
				Owner: sub.OwnerName,
			})

			if err != nil {
				return err
			}

			RepositoriesInfo[idx] = RepoInfo

			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	return RepositoriesInfo, nil
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
