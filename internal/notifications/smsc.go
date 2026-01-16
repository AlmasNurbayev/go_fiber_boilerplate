package notifications

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/config"
)

func SMSC_SendSms(cfg *config.Config, log1 *slog.Logger, phoneNumber string, message string) error {
	op := "notifications.SMSC_SendSms"
	log := log1.With(slog.String("op", op))

	host := cfg.SMSC_HOST + "rest/send/"
	body := strings.NewReader(`{"login":"` + cfg.SMSC_USER + `",
	"psw":"` + cfg.SMSC_PASSWORD + `",
	"phones":"` + phoneNumber + `",
	"mes":"` + message + `"}`)

	req, err := http.NewRequest("POST", host, body)
	if err != nil {
		log.Error("Api error:", slog.String("err", err.Error()))
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Api error:", slog.String("err", err.Error()))
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error("Error closing response body:", slog.String("err", err.Error()))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Error("Api error:", slog.String("err", resp.Status))
		return err
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Api error:", slog.String("err", err.Error()))
		return err
	}
	var response struct {
		Error string `json:"error"`
	}

	if err := json.Unmarshal(resBody, &response); err != nil {
		log.Error("Api error:", slog.String("err", err.Error()))
		return err
	}
	if response.Error != "" {
		log.Error("Api error:", slog.String("err", response.Error))
		return errors.New(response.Error)
	}

	return nil
}
