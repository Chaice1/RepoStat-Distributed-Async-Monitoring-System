package subscriber

import (
	"context"
	"log/slog"
	adapter_errors "repo-stat/api/internal/adapter/errors"
	"repo-stat/api/internal/domain"

	subscirberpb "repo-stat/proto/subscriber"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   subscirberpb.SubscriberClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   subscirberpb.NewSubscriberClient(conn),
	}, nil
}

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.pb.Ping(ctx, &subscirberpb.PingRequest{})
	if err != nil {
		c.log.Error("subscriber ping failed", "error", err)
		return domain.PingStatusDown
	}

	return domain.PingStatusUp
}

func (c *Client) CreateSubscription(ctx context.Context, repo string, owner string) error {
	_, err := c.pb.CreateSubscription(ctx, &subscirberpb.CreateSubscriptionRequest{
		RepoName:  repo,
		OwnerName: owner,
	})
	if err != nil {
		return adapter_errors.ErrorHandleFromGRPCToDomainWithLog(err, c.log, "CreateSubscription")
	}

	return nil

}

func (c *Client) GetSubscriptions(ctx context.Context) ([]*domain.Subscription, error) {
	resp, err := c.pb.GetSubscriptions(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, adapter_errors.ErrorHandleFromGRPCToDomainWithLog(err, c.log, "GetSubscriptions")
	}

	subscriptions := make([]*domain.Subscription, len(resp.GetSubscriptions()))

	for i, item := range resp.GetSubscriptions() {
		id, err := uuid.Parse(item.Id)
		if err != nil {
			return nil, adapter_errors.ErrorHandleFromGRPCToDomainWithLog(err, c.log, "parse string to uuid")
		}
		subscriptions[i] = &domain.Subscription{
			Id:        id,
			OwnerName: item.OwnerName,
			RepoName:  item.RepoName,
			CreatedAt: item.CreatedAt,
		}
	}

	return subscriptions, nil
}

func (c *Client) DeleteSubscription(ctx context.Context, repo string, owner string) error {
	_, err := c.pb.DeleteSubscription(ctx, &subscirberpb.DeleteSubscriptionRequest{
		RepoName:  repo,
		OwnerName: owner,
	})

	if err != nil {
		return adapter_errors.ErrorHandleFromGRPCToDomainWithLog(err, c.log, "DeleteSubscription")
	}
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
