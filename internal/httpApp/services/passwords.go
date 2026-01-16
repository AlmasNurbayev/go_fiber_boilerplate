package services

import (
	"context"
	"log/slog"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
)

func (s *AuthService) UpdatePassword(ctx context.Context, id int64, oldpassword string, newpassword string) error {
	op := "services.UpdatePassword"
	log := s.log.With(slog.String("op", op))

	userEntity, dbError := s.authStorage.GetUserById(ctx, id)
	if dbError != nil {
		log.Warn("error get user by id", slog.String("err", dbError.Message))
		return errorsApp.ErrUserNotFound.Error
	}

	isValidErr := lib.CheckPassword(userEntity.Password_hash.String, oldpassword)
	if isValidErr != nil {
		log.Warn("error verify password", slog.String("err", "password not match"))
		return errorsApp.ErrOldPasswordNotMatch.Error
	}

	err := s.authStorage.UpdatePassword(ctx, id, newpassword)
	if err != nil {
		log.Warn("error update password", slog.String("err", err.Message))
		return errorsApp.ErrInternalError.Error
	}

	log.Debug("password updated", slog.Int64("user_id", id))

	return nil
}
