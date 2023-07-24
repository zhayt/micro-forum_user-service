package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	AppPort    string `env:"APP_PORT" envDefault:"5001"`
	AppMode    string `env:"APP_MODE" envDefault:"dev"`
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBName     string `env:"DB_NAME"`
	DBPassword string `env:"DB_PASSWORD"`
	TZ         string `env:"TZ" envDefault:"Asia/Almaty"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config: %w", err)
	}

	return &cfg, nil
}

func PrepareENV() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
}
