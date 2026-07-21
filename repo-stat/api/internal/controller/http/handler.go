package http

import (
	"context"
	"log/slog"
	"net/http"
	"repo-stat/api/config"
	"repo-stat/api/internal/adapter/processor"
	api_redis "repo-stat/api/internal/adapter/redis"
	"repo-stat/api/internal/adapter/subscriber"
	"repo-stat/api/internal/usecase"
)

func NewHandler(ctx context.Context, log *slog.Logger, cfg config.Config) (http.Handler, error) {
	subscriberClient, err := subscriber.NewClient(cfg.Services.Subscriber, log)
	if err != nil {
		log.Error("cannot init subscriber adapter", "error", err)
		return nil, err
	}

	rRepo := api_redis.NewRedisRepo(cfg.Redis.Address, log, cfg.Cache.TTLSeconds)

	processorClient, err := processor.NewProcessorClient(cfg.Services.Processor, log, rRepo)
	if err != nil {
		log.Error("cannot init processor adapter", "error", err)
		return nil, err
	}

	pingUseCase := usecase.NewPing(subscriberClient)
	apiUsecase := usecase.NewUsecaseApiGateway(processorClient, subscriberClient)

	rL := NewRateLimiter(log, rRepo, cfg.RL.Burst, float64(cfg.RL.RPS))

	mux := http.NewServeMux()
	AddRoutes(mux, log, pingUseCase, apiUsecase)

	final_mux := rL.RateLimitMiddleware(mux)

	return final_mux, nil
}
