package processor_producer

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

type Repo interface {
	GetOutboxMessages(context.Context) ([]*processor_domain.RepoInfoTaskMessage, error)
	SetSentStatusOutboxMessage(context.Context, uuid.UUID) error
}

type producer struct {
	w   *kafka.Writer
	r   Repo
	log *slog.Logger
}

func NewProducer(address string, r Repo, log *slog.Logger) *producer {
	return &producer{
		w: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      []string{address},
			Topic:        "repo-tasks",
			Async:        false,
			RequiredAcks: -1,
		}),
		r:   r,
		log: log,
	}
}

func (p *producer) Relay(ctx context.Context) {

	jobs := make(chan *processor_domain.RepoInfoTaskMessage, 10)

	for i := 1; i <= 5; i++ {
		go func() {
			for {
				message := <-jobs

				dtoMessage := &processor_dto.RepoInfoTaskMessage{
					Id: message.Id,
					Payload: processor_dto.RepoInfoTaskMessagePayload{
						Repo:  message.Payload.Repo,
						Owner: message.Payload.Owner,
					},
					CreatedAt: message.CreatedAt.Format(time.RFC3339),
				}

				slice, _ := json.Marshal(&dtoMessage)

				var err error
				for try := 1; try <= 5; try++ {
					err = p.w.WriteMessages(ctx, kafka.Message{
						Key:   []byte(dtoMessage.Payload.Owner + "/" + dtoMessage.Payload.Repo),
						Value: slice,
					})

					if err != nil {
						p.log.Error("failed to write kafka message with key:%s", dtoMessage.Payload.Owner+"/"+dtoMessage.Payload.Repo, "error", err)
						time.Sleep(time.Second)
						continue
					}
					break
				}

				if err != nil {
					continue
				}

				err = p.r.SetSentStatusOutboxMessage(ctx, dtoMessage.Id)

				if err != nil {
					p.log.Error("failed to set sent status outbox message", "error", err)
					continue
				}

			}
		}()
	}

	for {

		messages, err := p.r.GetOutboxMessages(ctx)
		if err != nil || len(messages) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		for i := range messages {
			jobs <- messages[i]
		}

	}
}

func (p *producer) Close() error {
	return p.w.Close()
}
