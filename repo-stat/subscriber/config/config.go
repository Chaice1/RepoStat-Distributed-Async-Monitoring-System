package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-subscriber"`
}

type DataBase struct {
	Dsn  string `yaml:"dsn" env:"DB_DSN" env-default:"postgres://postgres:Ivbln173@localhost:5432/db?sslmode=disable"`
	Path string `yaml:"migration_path" env:"MIGRATION_PATH" env-default:"file://migrations"`
}

type Kafka struct {
	Address string `yaml:"address" env:"KAFKA_ADDR" env-default:"localhost:9092"`
}

type Services struct {
	API string `yaml:"api" env:"API_ADDRESS" env-default:"localhost:8080"`
}

type Config struct {
	App      App               `yaml:"app"`
	Services Services          `yaml:"services"`
	GRPC     grpcserver.Config `yaml:"grpc"`
	Logger   logger.Config     `yaml:"logger"`
	DB       DataBase          `yaml:"database"`
	Kafka    Kafka             `yaml:"kafka"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
