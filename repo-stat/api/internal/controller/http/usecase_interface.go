package http

import (
	"context"
	"repo-stat/api/internal/domain"
)

type ApiGatewayUsecase interface {
	GetInfoRepositories(context.Context) ([]*domain.RepoInfo, error)
	GetInfoRepo(context.Context, string, string) (*domain.RepoInfo, error)
	Ping(context.Context) domain.PingStatus
	DeleteSubscription(context.Context, string, string) error
	CreateSubscription(context.Context, string, string) error
	GetSubscriptions(context.Context) ([]*domain.Subscription, error)
}
