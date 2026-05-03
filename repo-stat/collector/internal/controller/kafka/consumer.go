package collector_consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	collectordomain "repo-stat/collector/internal/domain"
	collectorrespmodel "repo-stat/collector/internal/dto"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	SendRepoInfoMessage(context.Context, string, uuid.UUID, *collectorrespmodel.RepoInfoResMessage, error)
}

type UsecaseCollectorService interface {
	GetInfoRepo(context.Context, string, string) (*collectordomain.RepoInfo, error)
}

type consumer struct {
	p   Producer
	u   UsecaseCollectorService
	r   *kafka.Reader
	log *slog.Logger
}

func NewConsumer(address string, p Producer, u UsecaseCollectorService, log *slog.Logger) *consumer {
	return &consumer{
		p: p,
		u: u,
		r: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers: []string{address},
				Topic:   "repo-tasks",
				GroupID: "repo_info_1",
			},
		),
		log: log,
	}
}

func (c *consumer) Start(ctx context.Context, count_workers int) {
	for i := 1; i <= count_workers; i++ {
		go c.Consume(ctx, i)
	}
}

func (c *consumer) Consume(ctx context.Context, i int) {

	c.log.Debug("worker with id:%d started working", i)
	for {

		kafka_message, err := c.r.FetchMessage(ctx)

		if err != nil {
			continue
		}

		var message collectorrespmodel.RepoInfoMessage

		if err := json.Unmarshal(kafka_message.Value, &message); err != nil {
			c.log.Error("failed to parse json", "error", err)
			c.r.CommitMessages(ctx, kafka_message)
			continue
		}

		RepoInfo, err := c.u.GetInfoRepo(ctx, message.Payload.Owner, message.Payload.Repo)

		if err != nil {
			if errors.Is(err, collectordomain.ErrInternalError) {
				c.log.Error("failed to get info repo", "error", err)
				time.Sleep(1 * time.Second)
				continue
			}
			c.p.SendRepoInfoMessage(ctx, message.Payload.Owner+"/"+message.Payload.Repo, uuid.New(), nil, err)
			c.r.CommitMessages(ctx, kafka_message)
			continue
		}

		RepoInfoDto := collectorrespmodel.RepoInfoResMessage{
			Id:          uuid.New(),
			FullName:    RepoInfo.FullName,
			Description: RepoInfo.Description,
			Forks:       RepoInfo.Forks,
			Stargazers:  RepoInfo.Stargazers,
			CreatedAt:   RepoInfo.CreatedAt,
			Error:       "",
		}

		c.p.SendRepoInfoMessage(ctx, message.Payload.Owner+"/"+message.Payload.Repo, RepoInfoDto.Id, &RepoInfoDto, nil)
		c.r.CommitMessages(ctx, kafka_message)

	}
}

func (c *consumer) Close() error {
	return c.r.Close()
}
