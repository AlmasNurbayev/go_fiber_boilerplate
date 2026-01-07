package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	SERVICE_NAME string `env:"SERVICE_NAME,required"`

	POSTGRES_USER     string        `env:"POSTGRES_USER,required" json:"-"`
	POSTGRES_PASSWORD string        `env:"POSTGRES_PASSWORD,required" json:"-"`
	POSTGRES_DB       string        `env:"POSTGRES_DB,required"`
	POSTGRES_PORT     string        `env:"POSTGRES_PORT,required"`
	POSTGRES_TIMEOUT  time.Duration `env:"POSTGRES_TIMEOUT,required"`
	POSTGRES_HOST     string        `env:"POSTGRES_HOST,required"`

	REDIS_HOST       string `env:"REDIS_HOST,required"`
	REDIS_PORT       string `env:"REDIS_PORT,required"`
	REDIS_SESSION_DB int    `env:"REDIS_SESSION_DB,required"`

	AUTH_SECRET_KEY               string `env:"AUTH_SECRET_KEY,required"  json:"-"` // для шифрования в БД
	AUTH_ACCESS_TOKEN_EXP_MINUTES int    `env:"AUTH_ACCESS_TOKEN_EXP_HOURS"`
	AUTH_REFRESH_TOKEN_EXP_HOURS  int    `env:"AUTH_REFRESH_TOKEN_EXP_HOURS"`

	HTTP_PORT                   string        `env:"HTTP_PORT,required"`
	HTTP_TIMEOUT                time.Duration `env:"HTTP_TIMEOUT,required"`
	HTTP_PREFORK                bool          `env:"HTTP_PREFORK"`
	HTTP_CORS_ALLOW_ORIGINS     []string      `env:"HTTP_CORS_ALLOW_ORIGINS,required"`
	HTTP_CORS_ALLOW_CREDENTIALS bool          `env:"HTTP_CORS_ALLOW_CREDENTIALS,required"`
	HTTP_CORS_ALLOW_HEADERS     []string      `env:"HTTP_CORS_ALLOW_HEADERS,required"`

	PROMETHEUS_HTTP_PORT string `env:"PROMETHEUS_HTTP_PORT,required"`

	NATS_NAME            string `env:"NATS_NAME"`
	NATS_PORT            string `env:"NATS_PORT"`
	NATS_MONITORING_PORT string `env:"NATS_MONITORING_PORT"`
	NATS_STREAM_NAME     string `env:"NATS_STREAM_NAME"`

	LOG_ERROR_PATH string `env:"LOG_ERROR_PATH"`

	ENV string `env:"ENV,required"`
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
