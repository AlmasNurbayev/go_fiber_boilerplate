package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	POSTGRES_USER     string        `env:"POSTGRES_USER" json:"-"`
	POSTGRES_PASSWORD string        `env:"POSTGRES_PASSWORD" json:"-"`
	POSTGRES_DB       string        `env:"POSTGRES_DB"`
	POSTGRES_PORT     string        `env:"POSTGRES_PORT"`
	POSTGRES_TIMEOUT  time.Duration `env:"POSTGRES_TIMEOUT"`
	POSTGRES_HOST     string        `env:"POSTGRES_HOST"`
	POSTGRES_INT_PORT string        `env:"POSTGRES_INT_PORT"`

	AUTH_SECRET_KEY string `env:"SECRET_KEY"  json:"-"` // для шифрования в БД

	HTTP_PORT                   string        `env:"HTTP_PORT,required"`
	HTTP_TIMEOUT                time.Duration `env:"HTTP_TIMEOUT"`
	HTTP_PREFORK                bool          `env:"HTTP_PREFORK"`
	HTTP_CORS_ALLOW_ORIGINS     []string      `env:"HTTP_CORS_ALLOW_ORIGINS"`
	HTTP_CORS_ALLOW_CREDENTIALS bool          `env:"HTTP_CORS_ALLOW_CREDENTIALS"`
	HTTP_CORS_ALLOW_HEADERS     []string      `env:"HTTP_CORS_ALLOW_HEADERS"`

	PROMETHEUS_HTTP_PORT string `env:"PROMETHEUS_HTTP_PORT"`

	NATS_NAME            string `env:"NATS_NAME"`
	NATS_PORT            string `env:"NATS_PORT"`
	NATS_MONITORING_PORT string `env:"NATS_MONITORING_PORT"`
	NATS_STREAM_NAME     string `env:"NATS_STREAM_NAME"`

	LOG_ERROR_PATH string `env:"LOG_ERROR_PATH"`

	ENV string `env:"ENV"`
}

func Mustload(path string) *Config {
	cfg := &Config{}

	if path != "" {
		err := godotenv.Load(path)
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	err := env.Parse(cfg)
	if err != nil {
		log.Fatal("Error parse env: ", err)
	}

	// if cfg.SECRET_KEY != "" {
	// 	cfg.SECRET_BYTE = utils.DeriveKeyFromSecret(cfg.SECRET_KEY)
	// }

	if cfg.LOG_ERROR_PATH == "" {
		cfg.LOG_ERROR_PATH = "_volume_assets/error.log"
	}

	return cfg
}
