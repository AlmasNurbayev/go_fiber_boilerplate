package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/redis/go-redis/v9"
)

type OtpStorage struct {
	//Ctx context.Context
	RDB *redis.Client
	log *slog.Logger
}

func InitOtp(ctx context.Context, host string, port string, number int, log *slog.Logger) (*OtpStorage, error) {
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

	log.Info("Redis OTP storage initialized")

	return &OtpStorage{RDB: RDB, log: log}, nil
}

type OtpData struct {
	UserID    int64     `json:"user_id"`
	Otp       string    `json:"otp"`
	Type      string    `json:"type"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	ExpireAt  time.Time `json:"expire_at"`
}

func (c *OtpStorage) SaveOtp(
	ctx context.Context,
	data OtpData,
	ttlMinutes int,
) *errorsApp.DbError {

	op := "cache.OtpStorage.SaveOtp"
	log := c.log.With(slog.String("op", op))

	key := fmt.Sprintf("otp:%s:%s", data.Type, data.Address)
	//indexKey := fmt.Sprintf("otp:index:%s:%s", data.Type, data.Address)

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error("marshal otp data", slog.Any("err", err))
		return &errorsApp.DbError{
			Type:    "internal_error",
			Field:   "data",
			Message: "internal error marshal otp data",
			Error:   err,
		}
	}

	ttl := time.Duration(ttlMinutes) * time.Minute
	ok, err := c.RDB.SetNX(ctx, key, jsonData, ttl).Result()
	if err != nil {
		return &errorsApp.DbError{
			Type:    "internal_error",
			Field:   "data",
			Message: "internal error save otp",
			Error:   err,
		}
	}
	if !ok {
		return &errorsApp.DbError{
			Type:    "already_otp",
			Field:   "data",
			Message: "otp already sent, wait TTL",
			Error:   err,
		}
	}

	return nil
}

func (c *OtpStorage) DeleteOtp(
	ctx context.Context,
	address string,
	typeM string) *errorsApp.DbError {

	op := "cache.OtpStorage.DeleteOtp"
	log := c.log.With(slog.String("op", op))

	key := fmt.Sprintf("otp:%s:%s", typeM, address)

	_, err := c.RDB.Del(ctx, key).Result()
	if err != nil {
		log.Warn("delete otp", slog.Any("err", err))
		return &errorsApp.DbError{
			Type:    "internal_error",
			Field:   "data",
			Message: "internal error delete otp",
			Error:   err,
		}
	}
	return nil
}

func (c *OtpStorage) GetOtp(
	ctx context.Context,
	address string,
	typeM string) (OtpData, *errorsApp.DbError) {

	op := "cache.OtpStorage.GetOtp"
	log := c.log.With(slog.String("op", op))

	otpData := OtpData{}
	key := fmt.Sprintf("otp:%s:%s", typeM, address)

	data, err := c.RDB.Get(ctx, key).Result()
	if err != nil {
		log.Warn("get otp", slog.Any("err", err))
		return otpData, &errorsApp.DbError{
			Type:    "internal_error",
			Field:   "data",
			Message: "error or not found get otp",
			Error:   err,
		}
	}
	err2 := json.Unmarshal([]byte(data), &otpData)
	if err2 != nil {
		log.Error("error unmarshal otp data", slog.String("err", err2.Error()))
		return otpData, &errorsApp.DbError{
			Type:    "internal_error",
			Field:   "data",
			Message: "internal error unmarshal otp data",
			Error:   err2,
		}
	}

	return otpData, nil
}
