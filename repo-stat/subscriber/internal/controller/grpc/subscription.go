package grpc

import (
	"context"
	"log/slog"
	subscriberpb "repo-stat/proto/subscriber"
	grpc_errors "repo-stat/subscriber/internal/controller/errors"
	subscriber_domain "repo-stat/subscriber/internal/domain"
	"repo-stat/subscriber/internal/usecase"

	"google.golang.org/protobuf/types/known/emptypb"
)

type SubscriptionUsecase interface {
	DeleteSubscription(ctx context.Context, repo string, owner string) error
	CreateSubscription(ctx context.Context, repo string, owner string) error
	GetSubscriptions(ctx context.Context) ([]*subscriber_domain.Subscription, error)
}

type PingUsecase interface {
	Execute(context.Context) string
}

type subscriptionController struct {
	subscriberpb.UnimplementedSubscriberServer
	su   SubscriptionUsecase
	log  *slog.Logger
	ping *usecase.Ping
}

func NewSubscriptionController(su SubscriptionUsecase, log *slog.Logger, ping *usecase.Ping) *subscriptionController {
	return &subscriptionController{
		su:   su,
		log:  log,
		ping: ping,
	}
}

func (sc *subscriptionController) GetSubscriptions(ctx context.Context, req *emptypb.Empty) (*subscriberpb.GetSubscriptionsResponse, error) {

	subscriptions, err := sc.su.GetSubscriptions(ctx)

	if err != nil {
		return nil, grpc_errors.HandleErrorFromDomainToGRPC(err)
	}

	resp := make([]*subscriberpb.Subscription, len(subscriptions))

	for i := range subscriptions {
		resp[i] = &subscriberpb.Subscription{
			Id:        subscriptions[i].Id.String(),
			OwnerName: subscriptions[i].OwnerName,
			RepoName:  subscriptions[i].RepoName,
			CreatedAt: subscriptions[i].CreatedAt,
		}
	}
	return &subscriberpb.GetSubscriptionsResponse{
		Subscriptions: resp,
	}, nil
}

func (sc *subscriptionController) DeleteSubscription(ctx context.Context, req *subscriberpb.DeleteSubscriptionRequest) (*emptypb.Empty, error) {
	err := sc.su.DeleteSubscription(ctx, req.RepoName, req.OwnerName)

	if err != nil {
		return nil, grpc_errors.HandleErrorFromDomainToGRPC(err)
	}

	return &emptypb.Empty{}, nil

}

func (sc *subscriptionController) CreateSubscription(ctx context.Context, req *subscriberpb.CreateSubscriptionRequest) (*emptypb.Empty, error) {
	err := sc.su.CreateSubscription(ctx, req.RepoName, req.OwnerName)

	if err != nil {
		return nil, grpc_errors.HandleErrorFromDomainToGRPC(err)
	}

	return &emptypb.Empty{}, nil
}

func (sc *subscriptionController) Ping(ctx context.Context, _ *subscriberpb.PingRequest) (*subscriberpb.PingResponse, error) {
	sc.log.Debug("subscriberp ping request received")

	return &subscriberpb.PingResponse{
		Reply: sc.ping.Execute(ctx),
	}, nil
}
