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

	// проверяем существует ли пользователь
	if body.Type == "email" {
		user, err := s.authStorage.GetUserByEmail(ctx, body.Address)
		if err != nil {
			log.Warn("error get user by email", slog.String("err", err.Message))
			return errorsApp.ErrInternalError.Error
		}
		if user.Id == 0 {
			log.Warn("user not found by email", slog.String("address", body.Address))
			return errorsApp.ErrAuthentication.Error
		}
	}
	if body.Type == "phone" {
		user, err := s.authStorage.GetUserByPhoneNumber(ctx, body.Address)
		if err != nil {
			log.Warn("error get user by phone", slog.String("err", err.Message))
			return errorsApp.ErrInternalError.Error
		}
		if user.Id == 0 {
			log.Warn("user not found by phone", slog.String("address", body.Address))
			return errorsApp.ErrAuthentication.Error
		}
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
		go func() {
			err := notifications.SMSC_SendSms(s.cfg, s.log, body.Address, "Your verify code is: "+otp)
			if err != nil {
				log.Warn("error send verify code to user", slog.String("err", err.Error()))
			} else {
				log.Info("send verify code to user", slog.String("body", body.Address))
			}
		}()
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

func (s *AuthService) ConfirmVerify(ctx context.Context, body dto.AuthConfirmVerifyRequest) error {
	op := "services.ConfirmVerify"
	log := s.log.With(slog.String("op", op))

	if body.Type != "phone" && body.Type != "email" {
		log.Warn("invalid type", slog.String("type", body.Type))
		return errorsApp.ErrBadRequest.Error
	}
	if body.Address == "" {
		log.Warn("address is empty", slog.String("type", body.Type))
		return errorsApp.ErrBadRequest.Error
	}

	if body.Code == "" {
		log.Warn("code is empty", slog.String("type", body.Type))
		return errorsApp.ErrBadRequest.Error
	}

	otpData, err := s.otpStorage.GetOtp(ctx, body.Address, body.Type)
	if err != nil {
		log.Warn("error get otp", slog.String("err", err.Message))
		return errorsApp.ErrInternalError.Error
	}
	if otpData.Otp != body.Code {
		log.Warn("invalid otp", slog.String("otp", body.Code))
		return errorsApp.ErrAuthentication.Error
	}

	if body.Type == "phone" {
		// сначала ищем пользователя по телефону
		user, err := s.authStorage.GetUserByPhoneNumber(ctx, body.Address)
		if err != nil {
			log.Warn("error get user by phone", slog.String("err", err.Message))
			return errorsApp.ErrInternalError.Error
		}
		if user.Id == 0 {
			log.Warn("user not found by phone", slog.String("address", body.Address))
			return errorsApp.ErrAuthentication.Error
		}
		// если пользователь найден, обновляем время верификации
		err2 := s.authStorage.UpdateUserPhoneVerifyTimestamp(ctx, user.Id)
		if err2 != nil {
			log.Warn("error update user phone verify timestamp", slog.String("err", err.Message))
			return errorsApp.ErrInternalError.Error
		}
	}
	if body.Type == "email" {
		// сначала ищем пользователя по email
		user, err := s.authStorage.GetUserByEmail(ctx, body.Address)
		if err != nil {
			log.Warn("error get user by email", slog.String("err", err.Message))
			return errorsApp.ErrInternalError.Error
		}
		if user.Id == 0 {
			log.Warn("user not found by email", slog.String("address", body.Address))
			return errorsApp.ErrAuthentication.Error
		}
		// если пользователь найден, обновляем время верификации
		err2 := s.authStorage.UpdateUserEmailVerifyTimestamp(ctx, user.Id)
		if err2 != nil {
			log.Warn("error update user email verify timestamp", slog.String("err", err2.Message))
			return errorsApp.ErrInternalError.Error
		}
	}

	return nil
}
