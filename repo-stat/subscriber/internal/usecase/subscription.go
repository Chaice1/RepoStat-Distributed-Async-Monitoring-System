package usecase

import (
	"context"
	"errors"
	"fmt"
	subscriber_domain "repo-stat/subscriber/internal/domain"
	"time"
)

type SubscriberRepo interface {
	GetSubscriptions(context.Context) ([]*subscriber_domain.Subscription, error)
	CreateSubscription(context.Context, string, string) error
	GetSubscription(context.Context, string, string) (*subscriber_domain.Subscription, error)
	DeleteSubscriptionTransaction(context.Context, string, string) error
}

type CollectorClient interface {
	GetRepoInfo(context.Context, string, string) error
}

type subscriberUsecase struct {
	sr SubscriberRepo
	cc CollectorClient
}

func NewSubscriberUsecase(sr SubscriberRepo, cc CollectorClient) *subscriberUsecase {
	return &subscriberUsecase{
		sr: sr,
		cc: cc,
	}
}

func (su *subscriberUsecase) GetSubscriptions(ctx context.Context) ([]*subscriber_domain.Subscription, error) {
	return su.sr.GetSubscriptions(ctx)
}

func (su *subscriberUsecase) DeleteSubscription(ctx context.Context, repo string, owner string) error {

	err := su.sr.DeleteSubscriptionTransaction(ctx, repo, owner)
	if err != nil {
		return err
	}

	return nil
}

func (su *subscriberUsecase) CreateSubscription(ctx context.Context, repo string, owner string) error {

	_, err := su.sr.GetSubscription(ctx, repo, owner)
	start := time.Now()

	if err == nil {
		return subscriber_domain.ErrSubscriptionAlreadyExists
	}

	fmt.Printf("DB req:%v\n", time.Since(start))

	if !errors.Is(err, subscriber_domain.ErrNotFound) {
		return err
	}

	start = time.Now()

	err = su.cc.GetRepoInfo(ctx, owner, repo)
	if err != nil {
		return err
	}

	fmt.Printf("GithubAPI req:%v\n", time.Since(start))

	return su.sr.CreateSubscription(ctx, repo, owner)

}
