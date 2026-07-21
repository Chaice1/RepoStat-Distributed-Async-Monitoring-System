package subsciber_repository

import (
	"context"
	subscriber_repo_errors "repo-stat/subscriber/internal/adapter/repository/errors"
	generated_db "repo-stat/subscriber/internal/adapter/repository/gen"
	subscriber_domain "repo-stat/subscriber/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	q *generated_db.Queries
	p *pgxpool.Pool
}

func NewRepo(p *pgxpool.Pool) *repo {
	return &repo{
		q: generated_db.New(p),
		p: p,
	}
}

func (r *repo) CreateSubscription(ctx context.Context, repo string, owner string) error {
	err := r.q.CreateSubscription(ctx, generated_db.CreateSubscriptionParams{
		RepoName:  repo,
		OwnerName: owner,
	})

	if err != nil {
		return subscriber_domain.ErrInternalError
	}
	return nil
}

func (r *repo) DeleteSubscription(ctx context.Context, repo string, owner string) error {
	err := r.q.DeleteSubscription(ctx, generated_db.DeleteSubscriptionParams{
		RepoName:  repo,
		OwnerName: owner,
	})

	if err != nil {
		return subscriber_domain.ErrInternalError
	}
	return nil
}

func (r *repo) GetSubscriptions(ctx context.Context) ([]*subscriber_domain.Subscription, error) {
	responce, err := r.q.GetSubscriptions(ctx)
	if err != nil {
		return nil, subscriber_repo_errors.HandleErrorFromDBToDomain(err)
	}

	subscriptions := make([]*subscriber_domain.Subscription, len(responce))

	for i := range responce {
		subscriptions[i] = &subscriber_domain.Subscription{
			Id:        responce[i].ID,
			RepoName:  responce[i].RepoName,
			OwnerName: responce[i].OwnerName,
			CreatedAt: responce[i].CreatedAt.Time.Format(time.RFC3339),
		}
	}

	return subscriptions, nil
}

func (r *repo) GetSubscription(ctx context.Context, repo string, owner string) (*subscriber_domain.Subscription, error) {
	responce, err := r.q.GetSubscription(ctx, generated_db.GetSubscriptionParams{
		RepoName:  repo,
		OwnerName: owner,
	})

	if err != nil {
		return nil, subscriber_repo_errors.HandleErrorFromDBToDomain(err)
	}
	return &subscriber_domain.Subscription{
		Id:        responce.ID,
		OwnerName: responce.OwnerName,
		RepoName:  responce.RepoName,
		CreatedAt: responce.CreatedAt.Time.Format(time.RFC3339),
	}, nil
}

func (r *repo) GetOutboxMessage(ctx context.Context) (*subscriber_domain.RepoInfoTaskMessage, error) {

	responce, err := r.q.GetOutboxMessage(ctx)
	if err != nil {
		return nil, subscriber_repo_errors.HandleErrorFromDBToDomain(err)
	}

	return &subscriber_domain.RepoInfoTaskMessage{
		Id: responce.ID,
		Payload: subscriber_domain.RepoInfoTaskMessagePayload{
			Owner: responce.Owner,
			Repo:  responce.Repo,
		},
		CreatedAt: responce.CreatedAt.Time.Format(time.RFC3339),
	}, nil
}

func (r *repo) CreateOutboxMessage(ctx context.Context, repo string, owner string) error {

	err := r.q.CreateOutboxMessage(ctx, generated_db.CreateOutboxMessageParams{
		Owner: owner,
		Repo:  repo,
	})

	if err != nil {
		return subscriber_repo_errors.HandleErrorFromDBToDomain(err)
	}

	return nil
}

func (r *repo) SetSentStatusOutboxMessage(ctx context.Context, id uuid.UUID) error {
	err := r.q.SetSentStatusOutboxMessage(ctx, id)
	if err != nil {
		return subscriber_repo_errors.HandleErrorFromDBToDomain(err)
	}

	return nil
}

func (r *repo) DeleteSubscriptionTransaction(ctx context.Context, repo string, owner string) error {

	tx, err := r.p.Begin(ctx)

	if err != nil {
		return subscriber_repo_errors.HandleErrorFromDBToDomain(err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	qq := r.q.WithTx(tx)

	err = qq.DeleteSubscription(ctx, generated_db.DeleteSubscriptionParams{
		RepoName:  repo,
		OwnerName: owner,
	})

	if err != nil {
		return subscriber_repo_errors.HandleErrorFromDBToDomain(err)
	}

	err = qq.CreateOutboxMessage(ctx, generated_db.CreateOutboxMessageParams{
		Repo:  repo,
		Owner: owner,
	})

	if err != nil {
		return subscriber_repo_errors.HandleErrorFromDBToDomain(err)
	}

	err = tx.Commit(ctx)

	if err != nil {
		return subscriber_repo_errors.HandleErrorFromDBToDomain(err)
	}

	return nil
}
