package processor_repository

import (
	"context"
	"errors"
	"log/slog"
	processor_db_errors "repo-stat/processor/internal/adapter/repository/errors"
	generated_processor_db "repo-stat/processor/internal/adapter/repository/gen"
	processor_domain "repo-stat/processor/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	q   *generated_processor_db.Queries
	p   *pgxpool.Pool
	log *slog.Logger
}

func NewRepo(pool *pgxpool.Pool, log *slog.Logger) *repo {
	return &repo{
		q:   generated_processor_db.New(pool),
		p:   pool,
		log: log,
	}
}

func (r *repo) SetSentStatusOutboxMessage(ctx context.Context, id uuid.UUID) error {
	err := r.q.SetSentStatusOutboxMessage(ctx, id)

	if err != nil {
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "SetSentStatusOutboxMessage")
	}

	return nil
}

func (r *repo) GetOutboxMessages(ctx context.Context) ([]*processor_domain.RepoInfoTaskMessage, error) {

	messages, err := r.q.GetOutboxMessage(ctx)

	if err != nil {
		return nil, processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "GetOutboxMessage")
	}

	outboxMessages := make([]*processor_domain.RepoInfoTaskMessage, len(messages))

	for i, item := range messages {
		outboxMessages[i] = &processor_domain.RepoInfoTaskMessage{
			Id: item.ID,
			Payload: processor_domain.RepoInfoTaskMessagePayload{
				Repo:  item.Repo,
				Owner: item.Owner,
			},
			CreatedAt: item.CreatedAt.Time,
		}
	}

	return outboxMessages, nil
}

func (r *repo) CreateRepoTransaction(ctx context.Context, id uuid.UUID, payload []byte, repoinfo *processor_domain.RepoInfo) error {

	tx, err := r.p.Begin(ctx)
	if err != nil {
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "CreateRepoTransaction")
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()
	qq := r.q.WithTx(tx)

	err = qq.CreateInboxMessage(ctx, generated_processor_db.CreateInboxMessageParams{
		ID:      id,
		Payload: payload,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil
		}
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "CreateRepoTransaction")
	}

	err = qq.CreateOrUpdateRepoInfo(ctx, generated_processor_db.CreateOrUpdateRepoInfoParams{
		Lower:       repoinfo.FullName,
		Description: repoinfo.Description,
		Forks:       int32(repoinfo.Forks),
		Stargazers:  int32(repoinfo.Stargazers),
		CreatedAt:   pgtype.Timestamptz{Time: repoinfo.CreatedAt, Valid: true},
	})

	if err != nil {
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "CreateRepoTransaction")
	}
	err = tx.Commit(ctx)
	if err != nil {
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "CreateRepoTranscation")
	}
	return nil

}

func (r *repo) GetRepoInfo(ctx context.Context, fullname string) (*processor_domain.RepoInfo, error) {
	RepoInfo, err := r.q.GetRepositoryInfo(ctx, fullname)
	if err != nil {
		return nil, processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "GetRepoInfo")
	}

	return &processor_domain.RepoInfo{
		FullName:    RepoInfo.Fullname,
		Description: RepoInfo.Description,
		Forks:       uint64(RepoInfo.Forks),
		Stargazers:  uint64(RepoInfo.Stargazers),
		Status:      RepoInfo.Status,
		CreatedAt:   RepoInfo.CreatedAt.Time,
	}, nil
}

func (r *repo) DeleteRepoTransaction(ctx context.Context, fullname string, id uuid.UUID, payload []byte) error {

	tx, err := r.p.Begin(ctx)

	if err != nil {
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "DeleteRepo")
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	qq := r.q.WithTx(tx)

	err = qq.CreateInboxMessage(ctx, generated_processor_db.CreateInboxMessageParams{
		ID:      id,
		Payload: payload,
	})

	if err != nil {
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "DeleteRepo")
	}

	err = qq.DeleteRepo(ctx, fullname)

	if err != nil {
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "DeleteRepo")
	}

	err = tx.Commit(ctx)

	if err != nil {
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "DeleteRepo")
	}

	return nil
}

func (r *repo) CreateFetchingTaskTransaction(ctx context.Context, repo string, owner string) error {
	tx, err := r.p.Begin(ctx)

	if err != nil {
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "CreateFetchingTask")
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	qq := r.q.WithTx(tx)

	err = qq.CreateOutboxMessage(ctx, generated_processor_db.CreateOutboxMessageParams{
		Repo:  repo,
		Owner: owner,
	})

	if err != nil {
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "CreateFetchingTask")
	}

	err = qq.CreateFetchingTask(ctx, owner+"/"+repo)

	if err != nil {
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "CreateFetchingTask")
	}

	err = tx.Commit(ctx)

	if err != nil {
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "CreateFetchingTask")
	}

	return nil
}

func (r *repo) SetErrorStatusRepo(ctx context.Context, fullname string, payload []byte, id uuid.UUID) error {
	tx, err := r.p.Begin(ctx)

	if err != nil {
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "SetErrorStatusRepo")
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	qq := r.q.WithTx(tx)

	err = qq.CreateInboxMessage(ctx, generated_processor_db.CreateInboxMessageParams{
		ID:      id,
		Payload: payload,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil
		}
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "SetErrorStatusRepo")
	}

	err = qq.SetErrorStatusRepo(ctx, fullname)

	if err != nil {
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "SetErrorStatusRepo")
	}

	err = tx.Commit(ctx)

	if err != nil {
		return processor_db_errors.ErrorHandleFromDBToDomain(err, r.log, "SetErrorStatusRepo")
	}

	return nil
}
