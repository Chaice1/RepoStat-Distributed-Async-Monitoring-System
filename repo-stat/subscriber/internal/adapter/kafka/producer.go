package subscriber_producer

import (
	"context"
	"encoding/json"
	"log/slog"
	subscriber_domain "repo-stat/subscriber/internal/domain"
	subscriber_dto "repo-stat/subscriber/internal/dto"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Repository interface {
	GetOutboxMessage(context.Context) (*subscriber_domain.RepoInfoTaskMessage, error)
	SetSentStatusOutboxMessage(context.Context, uuid.UUID) error
}
type producer struct {
	w   *kafka.Writer
	r   Repository
	log *slog.Logger
}

func NewProducer(address string, log *slog.Logger, r Repository) *producer {
	return &producer{
		w: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      []string{address},
			Topic:        "repo-results",
			Async:        false,
			RequiredAcks: -1,
		}),
		log: log,
		r:   r,
	}
}

func (p *producer) SendMessage(ctx context.Context) {

	for {
		message, err := p.r.GetOutboxMessage(ctx)

		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		DelSubscrMessage := &subscriber_dto.DeleteSubscriptionMessage{
			Id:        uuid.New(),
			FullName:  message.Payload.Owner + "/" + message.Payload.Repo,
			CreatedAt: message.CreatedAt,
			Error:     "del_subscription",
		}

		DelSubscrSlice, err := json.Marshal(DelSubscrMessage)

		if err != nil {
			p.log.Error("failed to serialize message", "error", err)
			continue
		}

		p.log.Debug("message", DelSubscrMessage)

		err = p.WriteDelSubscriptionMessage(ctx, DelSubscrMessage.Id, DelSubscrSlice)

		if err != nil {
			p.log.Error("failed to send message", "error", err)
			continue
		}

		err = p.r.SetSentStatusOutboxMessage(ctx, message.Id)
		p.log.Error("message", err)
		if err != nil {
			p.log.Error("failed to set sent status", "error", err)
			continue
		}
	}
}
func (p *producer) WriteDelSubscriptionMessage(ctx context.Context, id uuid.UUID, payload []byte) error {

	for i := 1; i <= 5; i++ {
		err := p.w.WriteMessages(ctx, kafka.Message{
			Key:   []byte(id.String()),
			Value: payload,
		})
		if err != nil {
			p.log.Error("failed to write kafka message with key", id)
			time.Sleep(50 * time.Millisecond)
			continue
		} else {
			return nil
		}
	}

	return subscriber_domain.ErrInternalError
}
