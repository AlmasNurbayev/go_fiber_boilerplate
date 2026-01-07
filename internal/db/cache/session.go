package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type SessionData struct {
	Jti       string    `json:"jti"`
	UserID    int64     `json:"user_id"`
	RoleID    int64     `json:"role_id"`
	UserAgent string    `json:"user_agent"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *SessionStorage) SaveSession(ctx context.Context, jti string, data SessionData, ttlHours int) error {
	data.CreatedAt = time.Now()
	data.Jti = jti
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	pipe := c.RDB.TxPipeline()
	pipe.Set(ctx, jti, jsonData, time.Duration(ttlHours)*time.Hour).Err()
	pipe.SAdd(ctx, strconv.FormatInt(data.UserID, 10), jti).Err()
	_, err = pipe.Exec(ctx)

	return err
}

func (c *SessionStorage) GetSessionByJti(ctx context.Context, jti string) (SessionData, error) {
	val, err := c.RDB.Get(ctx, jti).Result()
	if err != nil {
		return SessionData{}, err
	}

	var data SessionData
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return SessionData{}, err
	}
	return data, nil
}

func (c *SessionStorage) GetSessionsByUserId(ctx context.Context, userId int64) ([]SessionData, error) {

	indexKey := fmt.Sprintf("%d", userId)

	jtis, err := c.RDB.SMembers(ctx, indexKey).Result()
	if err != nil {
		return nil, err
	}

	var sessions []SessionData

	for _, jti := range jtis {
		sessionKey := jti
		raw, err := c.RDB.Get(ctx, sessionKey).Bytes()
		if err != nil {
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

func (c *SessionStorage) DeleteSessionByJti(ctx context.Context, jti string) error {

	data, err := c.RDB.Get(ctx, jti).Bytes()
	if err != nil {
		return err
	}

	var sessionData SessionData
	err = json.Unmarshal(data, &sessionData)
	if err != nil {
		return err
	}
	userIndexKey := fmt.Sprintf("%d", sessionData.UserID)

	pipe := c.RDB.TxPipeline()
	pipe.Del(ctx, jti)
	pipe.SRem(ctx, userIndexKey, jti)

	_, err = pipe.Exec(ctx)
	return err
}
