package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	subscriberpb "repo-stat/proto/subscriber"
	"repo-stat/subscriber/config"
	subscriber_grpc_client "repo-stat/subscriber/internal/adapter/client/grpc"
	subscriber_producer "repo-stat/subscriber/internal/adapter/kafka"
	subsciber_repository "repo-stat/subscriber/internal/adapter/repository/postgres"
	grpccontroller "repo-stat/subscriber/internal/controller/grpc"
	"repo-stat/subscriber/internal/usecase"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func runMigrations(dsn string, path string) error {
	m, err := migrate.New(path, dsn)

	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil

}
func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting subscriber server...")
	log.Debug("debug messages are enabled")

	err := runMigrations(cfg.DB.Dsn, cfg.DB.Path)
	if err != nil {
		return err
	}

	pool, err := pgxpool.New(ctx, cfg.DB.Dsn)
	if err != nil {
		return err
	}

	repo := subsciber_repository.NewRepo(pool)

	producer := subscriber_producer.NewProducer(cfg.Kafka.Address, log, repo)

	gitHubApiClient := subscriber_grpc_client.NewGitHubApiClient()

	subscriberUsecase := usecase.NewSubscriberUsecase(repo, gitHubApiClient)
	pingUseCase := usecase.NewPing()
	subscriberController := grpccontroller.NewSubscriptionController(subscriberUsecase, log, pingUseCase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	go producer.SendMessage(ctx)

	subscriberpb.RegisterSubscriberServer(srv.GRPC(), subscriberController)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

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
