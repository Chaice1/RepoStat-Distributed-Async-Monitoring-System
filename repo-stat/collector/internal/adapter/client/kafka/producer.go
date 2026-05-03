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

func (p *producer) sendMessage(ctx context.Context, key []byte, Value []byte, Topic string) {

	for try := 1; try <= 5; try++ {

		err := p.w.WriteMessages(ctx, kafka.Message{
			Key:   key,
			Value: Value,
			Topic: Topic,
		})
		if err != nil {
			p.log.Error("failed to write kafka message with key:%s", string(key), "error", err)
			time.Sleep(time.Second)
			continue
		}
		return

	}
}

func (p *producer) SendRepoInfoMessage(ctx context.Context, fullname string, id uuid.UUID, RepoInfo *collectorrespmodel.RepoInfoResMessage, Err error) {

	if RepoInfo == nil {

		RepoInfoByteSlice, _ := json.Marshal(&collectorrespmodel.RepoInfoResMessage{
			FullName: fullname,
			Error:    Err.Error(),
		})

		p.sendMessage(ctx, []byte(id.String()), RepoInfoByteSlice, "repo-results")
		return
	}

	SliceByte, _ := json.Marshal(RepoInfo)

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
