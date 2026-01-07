package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
)

type SessionData struct {
	Jti       string    `json:"jti"`
	UserID    int64     `json:"user_id"`
	RoleID    int64     `json:"role_id"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *SessionStorage) SaveSession(ctx context.Context, jti string, data SessionData, ttlHours int) *errorsApp.DbError {
	op := "cache.SessionStorage.SaveSession"
	log := c.log.With(slog.String("op", op))

	data.CreatedAt = time.Now()
	data.Jti = jti
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error("error marshal session data", slog.String("err", err.Error()))
		return &errorsApp.DbError{
			Type:    "internal_error",
			Field:   "data",
			Message: "internal error marshal session data",
			Error:   err,
		}
	}

	pipe := c.RDB.TxPipeline()
	pipe.Set(ctx, jti, jsonData, time.Duration(ttlHours)*time.Hour).Err()
	pipe.SAdd(ctx, strconv.FormatInt(data.UserID, 10), jti).Err()
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Error("error save session", slog.String("err", err.Error()))
		return &errorsApp.DbError{
			Type:    "internal_error",
			Field:   "data",
			Message: "internal error save session",
			Error:   err,
		}
	}

	return nil
}

func (c *SessionStorage) GetSessionByJti(ctx context.Context, jti string) (SessionData, *errorsApp.DbError) {
	op := "cache.SessionStorage.GetSessionByJti"
	log := c.log.With(slog.String("op", op))

	val, err := c.RDB.Get(ctx, jti).Result()
	if err != nil {
		log.Error("error get session by jti", slog.String("err", err.Error()))
		return SessionData{}, &errorsApp.DbError{
			Type:    "not_found",
			Field:   "jti",
			Message: "session not found",
			Error:   err,
		}
	}

	var data SessionData
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		log.Error("error unmarshal session data", slog.String("err", err.Error()))
		return SessionData{}, &errorsApp.DbError{
			Type:    "internal_error",
			Field:   "data",
			Message: "internal error unmarshal session data",
			Error:   err,
		}
	}
	return data, nil
}

func (c *SessionStorage) GetSessionsByUserId(ctx context.Context, userId int64) ([]SessionData, *errorsApp.DbError) {

	op := "cache.SessionStorage.GetSessionByJti"
	log := c.log.With(slog.String("op", op))

	indexKey := fmt.Sprintf("%d", userId)

	jtis, err := c.RDB.SMembers(ctx, indexKey).Result()
	if err != nil {
		log.Error("error get sessions by user id", slog.String("err", err.Error()))
		return nil, &errorsApp.DbError{
			Type:    "not_found",
			Field:   "index",
			Message: "session index not found",
			Error:   err,
		}
	}

	var sessions []SessionData

	for _, jti := range jtis {
		sessionKey := jti
		raw, err := c.RDB.Get(ctx, sessionKey).Bytes()
		if err != nil {
			log.Error("error get session by jti", slog.String("err", err.Error()))
			// сессия могла протухнуть - чистим индекс
			_ = c.RDB.SRem(ctx, indexKey, jti).Err()
			continue
		}
		var s SessionData
		if err := json.Unmarshal(raw, &s); err != nil {
			continue
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (c *SessionStorage) DeleteSessionByJti(ctx context.Context, jti string) *errorsApp.DbError {

	op := "cache.SessionStorage.DeleteSessionByJti"
	log := c.log.With(slog.String("op", op))

	data, err := c.RDB.Get(ctx, jti).Bytes()
	if err != nil {
		log.Error("error get session by jti", slog.String("err", err.Error()))
		return &errorsApp.DbError{
			Type:    "not_found",
			Field:   "jti",
			Message: "session not found",
			Error:   err,
		}
	}

	var sessionData SessionData
	err = json.Unmarshal(data, &sessionData)
	if err != nil {
		log.Error("error unmarshal session data", slog.String("err", err.Error()))
		return &errorsApp.DbError{
			Type:    "internal_error",
			Field:   "data",
			Message: "internal error unmarshal session data",
			Error:   err,
		}
	}
	userIndexKey := fmt.Sprintf("%d", sessionData.UserID)

	pipe := c.RDB.TxPipeline()
	pipe.Del(ctx, jti)
	pipe.SRem(ctx, userIndexKey, jti)

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Error("error delete session by jti", slog.String("err", err.Error()))
		return &errorsApp.DbError{
			Type:    "internal_error",
			Field:   "data",
			Message: "internal error delete session by jti",
			Error:   err,
		}
	}
	return nil
}
