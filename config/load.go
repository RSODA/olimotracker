package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
}

type AppConfig struct {
	Host string `env:"APP_HOST" env-required:"true"`
	Port int    `env:"APP_PORT" env-required:"true"`
}

type PostgresConfig struct {
	Host     string `env:"PG_HOST" env-required:"true"`
	Port     int    `env:"PG_PORT" env-required:"true"`
	User     string `env:"PG_USER" env-required:"true"`
	Password string `env:"PG_PASSWORD" env-required:"true"`
	DB       string `env:"PG_DB" env-required:"true"`
}

func Load() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disbale",
		p.Host, p.Port, p.User, p.Password, p.DB,
	)
}
