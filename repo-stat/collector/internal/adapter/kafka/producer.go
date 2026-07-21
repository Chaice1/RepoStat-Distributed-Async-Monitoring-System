package collector_producer

import (
	"context"
	"encoding/json"
	"log/slog"
	collectordomain "repo-stat/collector/internal/domain"
	collectorrespmodel "repo-stat/collector/internal/dto"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type SubscriberClient interface {
	GetSubscriptions(context.Context) ([]*collectordomain.Subscription, error)
}

type producer struct {
	w   *kafka.Writer
	sc  SubscriberClient
	log *slog.Logger
}

func NewProducer(address string, sc SubscriberClient, log *slog.Logger) *producer {
	return &producer{
		w: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      []string{address},
			Async:        false,
			RequiredAcks: -1,
		}),
		sc:  sc,
		log: log,
	}
}

func (p *producer) sendMessage(ctx context.Context, key []byte, value []byte, topic string) {

	for try := 1; try <= 5; try++ {

		err := p.w.WriteMessages(ctx, kafka.Message{
			Key:   key,
			Value: value,
			Topic: topic,
		})
		if err != nil {
			p.log.Error("failed to write kafka message with key", "key", key, "error", err)
			time.Sleep(time.Second)
			continue
		}
		return

	}
}

func (p *producer) SendRepoInfoMessage(ctx context.Context, fullname string, id uuid.UUID, repoInfo *collectorrespmodel.RepoInfoResMessage, err error) {

	if repoInfo == nil {

		RepoInfoByteSlice, _ := json.Marshal(&collectorrespmodel.RepoInfoResMessage{
			FullName: fullname,
			Error:    err.Error(),
		})

		p.sendMessage(ctx, []byte(id.String()), RepoInfoByteSlice, "repo-results")
		return
	}

	SliceByte, _ := json.Marshal(repoInfo)

	p.sendMessage(ctx, []byte(id.String()), SliceByte, "repo-results")
}

func (p *producer) SendSubscriptionMessage(ctx context.Context) {

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

			subscriptions, err := p.sc.GetSubscriptions(ctx)

			if err != nil {
				p.log.Error("failed to get subscriptions", "error", err)
				continue
			}

			for _, item := range subscriptions {

				subscriptionDTO := collectorrespmodel.RepoInfoMessage{
					Id: item.Id,
					Payload: collectorrespmodel.RepoInfoMessagePayload{
						Repo:  item.RepoName,
						Owner: item.OwnerName,
					},
					CreatedAt: item.CreatedAt,
				}

				subscriptionByteSlice, _ := json.Marshal(&subscriptionDTO)

				p.sendMessage(ctx, []byte(subscriptionDTO.Payload.Owner+"/"+subscriptionDTO.Payload.Repo), subscriptionByteSlice, "repo-tasks")

			}

		}
	}

}

func (p *producer) Close() error {
	return p.w.Close()
}
