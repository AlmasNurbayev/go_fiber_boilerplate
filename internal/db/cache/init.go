package cache

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type SessionStorage struct {
	//Ctx context.Context
	RDB *redis.Client
	log *slog.Logger
}

func InitSession(ctx context.Context, host string, port string, number int, log *slog.Logger) (*SessionStorage, error) {
	RDB := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       number, // Используем стандартную БД
	})

	// Проверка соединения
	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Error(fmt.Sprintf("Failed to connect to Redis: %v", err))
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Info("Redis session storage initialized")

	return &SessionStorage{RDB: RDB, log: log}, nil
}
