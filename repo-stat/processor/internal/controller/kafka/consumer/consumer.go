package processor_consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	processor_domain "repo-stat/processor/internal/domain"
	processor_dto "repo-stat/processor/internal/dto"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Repository interface {
	CreateRepoTransaction(context.Context, uuid.UUID, []byte, *processor_domain.RepoInfo) error
	DeleteRepoTransaction(context.Context, string, uuid.UUID, []byte) error
	SetErrorStatusRepo(context.Context, string, []byte, uuid.UUID) error
}
type consumer struct {
	r   *kafka.Reader
	rep Repository
	log *slog.Logger
}

func NewConsumer(address string, rep Repository, log *slog.Logger) *consumer {
	return &consumer{
		r: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{address},
			Topic:   "repo-results",
			GroupID: "processor_consumer_1",
		}),
		rep: rep,
		log: log,
	}
}
func (c *consumer) Start(ctx context.Context, workers int) {
	for i := 1; i <= workers; i++ {
		c.log.Debug("worker %d started working", i)
		go c.ConsumeMessage(ctx)
	}
}
func (c *consumer) ConsumeMessage(ctx context.Context) {

	for {

		message, err := c.r.FetchMessage(ctx)

		if err != nil {
			c.log.Error("failed to fetch message", "error", err)
			continue
		}

		var RepoInfoResMes processor_dto.RepoInfoResMessage

		err = json.Unmarshal(message.Value, &RepoInfoResMes)

		if err != nil {
			c.log.Error("failed to parse json", "error", err)
			c.r.CommitMessages(ctx, message)
			continue
		}

		if RepoInfoResMes.Error == "del_subscription" {
			err = c.rep.DeleteRepoTransaction(ctx, RepoInfoResMes.FullName, RepoInfoResMes.Id, message.Value)

			if err != nil {
				c.log.Error("failed to commit delete repo transaction", "error", err)
				continue
			}

			c.r.CommitMessages(ctx, message)
			continue
		}

		if RepoInfoResMes.Error != "" {
			c.log.Error("message has error", "error", RepoInfoResMes.Error)
			err = c.rep.SetErrorStatusRepo(ctx, RepoInfoResMes.FullName, message.Value, RepoInfoResMes.Id)
			if err != nil {
				c.log.Error("failed to set error status", "error", err)
				continue
			}
			c.r.CommitMessages(ctx, message)
			continue
		}

		time, err := time.Parse(time.RFC3339, RepoInfoResMes.CreatedAt)

		if err != nil {
			c.log.Error("failed to parse time", "error", err)
			c.r.CommitMessages(ctx, message)
			continue
		}

		RepoInfo := &processor_domain.RepoInfo{
			FullName:    RepoInfoResMes.FullName,
			Description: RepoInfoResMes.Description,
			Forks:       RepoInfoResMes.Forks,
			Stargazers:  RepoInfoResMes.Stargazers,
			CreatedAt:   time,
		}

		err = c.rep.CreateRepoTransaction(ctx, RepoInfoResMes.Id, message.Value, RepoInfo)

		if err != nil {
			c.log.Error("failed to commit repo transaction")
			continue
		}

		c.r.CommitMessages(ctx, message)
	}
}

func (c *consumer) Close() error {
	return c.r.Close()
}
