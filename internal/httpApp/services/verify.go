package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/db/cache"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp/dto"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/notifications"
)

func (s *AuthService) SendVerify(ctx context.Context, body dto.AuthSendVerifyRequest) error {
	op := "services.SendVerify"
	log := s.log.With(slog.String("op", op))

	if body.Type != "phone" && body.Type != "email" {
		log.Warn("invalid type", slog.String("type", body.Type))
		return errorsApp.ErrBadRequest.Error
	}
	if body.Address == "" {
		log.Warn("address is empty", slog.String("type", body.Type))
		return errorsApp.ErrBadRequest.Error
	}

	errDeleteOtp := s.otpStorage.DeleteOtp(ctx, body.Address, body.Type)
	if errDeleteOtp != nil {
		log.Warn("error delete otp", slog.String("err", errDeleteOtp.Message))
	}

	otp := lib.GenerateOTP()
	otpData := cache.OtpData{
		Otp:       otp,
		Type:      body.Type,
		Address:   body.Address,
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(time.Duration(s.cfg.AUTH_OTP_TTL_MINUTES) * time.Minute),
	}

	err := s.otpStorage.SaveOtp(ctx, otpData, s.cfg.AUTH_OTP_TTL_MINUTES)
	if err != nil {
		log.Warn("error save otp", slog.String("err", err.Message))
		if err.Type == "already_otp" {
			return errorsApp.ErrAlreadyOtp.Error
		}
		return errorsApp.ErrInternalError.Error
	}

	if body.Type == "phone" {
		log.Info("send verify code to user", slog.String("body", body.Address))
	}
	if body.Type == "email" {
		// Отправляем email асинхронно, чтобы ошибки не блокировали основной поток
		go func() {
			err := notifications.SendMail(s.cfg, body.Address, "Verify code for "+s.cfg.SERVICE_NAME, "Your verify code is: "+otp)
			if err != nil {
				log.Warn("error send verify code to user", slog.String("err", err.Error()))
			} else {
				log.Info("send verify code to user", slog.String("body", body.Address))
			}
		}()
	}

	return nil
}
