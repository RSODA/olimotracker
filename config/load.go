package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Http        HttpConfig
	Postgres    PostgresConfig
	LoggerLevel string `env:"LOGGER_LEVEL" env-required:"true"`
	JWT         JWTConfig
}

type HttpConfig struct {
	Host string `env:"HTTP_HOST" env-required:"true"`
	Port int    `env:"HTTP_PORT" env-required:"true"`
}

type JWTConfig struct {
	Secret  string `env:"JWT_SECRET" env-required:"true"`
	JWT_TTL int    `env:"JWT_TTL" env-required:"true"`
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
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.DB,
	)
}

func (h HttpConfig) Addr() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}
