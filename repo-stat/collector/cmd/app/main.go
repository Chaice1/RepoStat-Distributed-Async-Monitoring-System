package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	collector_config "repo-stat/collector/config"
	collector_grpc_client "repo-stat/collector/internal/adapter/client/grpc"
	collectorclient "repo-stat/collector/internal/adapter/client/http"
	collector_producer "repo-stat/collector/internal/adapter/kafka"
	collector_consumer "repo-stat/collector/internal/controller/kafka"

	collectorusecase "repo-stat/collector/internal/usecase"
	"repo-stat/platform/logger"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := collector_config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting collector server...")
	log.Debug("debug messages are enabled")

	gitHubClient := collectorclient.NewGitHubApiClient()

	subscriberClient, err := collector_grpc_client.NewSubscriberClient(cfg.Services.Subscriber)

	if err != nil {
		return fmt.Errorf("create grpc client %w", err)
	}

	producer := collector_producer.NewProducer(cfg.Kafka.Address, subscriberClient, log)
	defer func() {
		_ = producer.Close()
	}()

	go producer.SendSubscriptionMessage(ctx)

	collectorUseCase := collectorusecase.NewCollectorService(gitHubClient)

	consumer := collector_consumer.NewConsumer(cfg.Kafka.Address, producer, collectorUseCase, log)

	defer func() {
		_ = consumer.Close()
	}()
	go consumer.Start(ctx, 5)

	<-ctx.Done()

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("launching server error: %s\n", err)
		}
		cancel()
		os.Exit(1)
	}
}
