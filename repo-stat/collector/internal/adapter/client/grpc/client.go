package collector_grpc_client

import (
	"context"
	collectordomain "repo-stat/collector/internal/domain"
	subscriberv1 "repo-stat/proto/subscriber"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type subscriberClient struct {
	sc   subscriberv1.SubscriberClient
	conn *grpc.ClientConn
}

func NewSubscriberClient(address string) (*subscriberClient, error) {

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &subscriberClient{sc: subscriberv1.NewSubscriberClient(conn), conn: conn}, nil
}

func (sc *subscriberClient) GetSubscriptions(ctx context.Context) ([]*collectordomain.Subscription, error) {
	resp, err := sc.sc.GetSubscriptions(ctx, &emptypb.Empty{})

	if err != nil {
		status, ok := status.FromError(err)

		if !ok {
			return nil, collectordomain.ErrInternalError
		}

		switch status.Code() {
		case codes.NotFound:
			return nil, collectordomain.ErrNotFound
		default:
			return nil, collectordomain.ErrInternalError
		}

	}

	subscriptions := make([]*collectordomain.Subscription, len(resp.GetSubscriptions()))

	for i, item := range resp.GetSubscriptions() {
		id, err := uuid.Parse(item.Id)
		if err != nil {
			return nil, collectordomain.ErrInternalError
		}
		subscriptions[i] = &collectordomain.Subscription{
			Id:        id,
			RepoName:  item.RepoName,
			OwnerName: item.OwnerName,
			CreatedAt: item.CreatedAt,
		}
	}

	return subscriptions, nil
}
