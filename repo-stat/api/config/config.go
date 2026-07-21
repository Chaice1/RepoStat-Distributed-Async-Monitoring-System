package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/httpserver"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-api"`
}

type Services struct {
	Subscriber string `yaml:"subscriber" env:"SUBSCRIBER_ADDRESS" env-default:"localhost:8081"`
	Processor  string `yaml:"processor" env:"PROCESSOR_ADDRESS" env-default:"localhost:8082"`
}

type RedisDB struct {
	Address string `yaml:"address" env:"REDIS_ADDR" env-default:"localhost:6379"`
}

type Cache struct {
	TTLSeconds int `yaml:"ttl_seconds" env:"TTL_SECONDS" env-default:"60"`
}

type RateLimit struct {
	Burst int `yaml:"burst" env:"BURST" env-default:"10"`
	RPS   int `yaml:"requests_per_second" env:"RPS" env-default:"5"`
}

type Config struct {
	App      App               `yaml:"app"`
	Services Services          `yaml:"services"`
	Redis    RedisDB           `yaml:"redis"`
	HTTP     httpserver.Config `yaml:"http"`
	Cache    Cache             `yaml:"cache"`
	RL       RateLimit         `yaml:"rate_limit"`
	Logger   logger.Config     `yaml:"logger"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
