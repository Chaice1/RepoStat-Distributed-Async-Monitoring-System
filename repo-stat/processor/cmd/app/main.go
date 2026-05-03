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
	processor_config "repo-stat/processor/config"
	processor_producer "repo-stat/processor/internal/adapter/kafka/producer"
	processor_repository "repo-stat/processor/internal/adapter/repository/postgres"
	processor_controller "repo-stat/processor/internal/controller/grpc"
	processor_consumer "repo-stat/processor/internal/controller/kafka/consumer"
	processor_usecase "repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"

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

	cfg := processor_config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting processor server...")
	log.Debug("debug messages are enabled")

	err := runMigrations(cfg.DB.Dsn, cfg.DB.Path)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	pool, err := pgxpool.New(ctx, cfg.DB.Dsn)

	if err != nil {
		return fmt.Errorf("failed to create pool: %w", err)
	}

	DB := processor_repository.NewRepo(pool, log)

	producer := processor_producer.NewProducer(cfg.Kafka.Address, DB, log)

	go producer.Relay(ctx)

	processorUseCase := processor_usecase.NewProcessorService(DB)

	processorHandler := processor_controller.NewProcessorService(processorUseCase, log)

	consumer := processor_consumer.NewConsumer(cfg.Kafka.Address, DB, log)

	go consumer.Start(ctx, 5)
	srv, err := grpcserver.New(cfg.GRPC.Address)

	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	processorpb.RegisterProcessorServer(srv.GRPC(), processorHandler)

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
