package processor_config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type APP struct {
	Processor string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-processor"`
}

type DataBase struct {
	Dsn  string `yaml:"dsn" env:"DB_DSN" env-default:"postgres://postgres:Ivbln173@localhost:5432/db?sslmode=disable"`
	Path string `yaml:"migration_path" env:"MIGRATION_PATH" env-default:"file://processor/migrations"`
}

type Kafka struct {
	Address string `yaml:"address" env:"KAFKA_ADDR" env-default:"localhost:9092"`
}

type Config struct {
	App    APP               `yaml:"app"`
	Logger logger.Config     `yaml:"logger"`
	DB     DataBase          `yaml:"database"`
	Kafka  Kafka             `yaml:"kafka"`
	GRPC   grpcserver.Config `yaml:"grpc"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
