package services

import (
	"context"
	"log/slog"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/config"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp/dto"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/models"
)

type UserService struct {
	log         *slog.Logger
	userStorage userStorage
	cfg         *config.Config
}

type userStorage interface {
	GetUserById(ctx context.Context, id int64) (models.UserEntity, *errorsApp.DbError)
	GetUserByNameStorage(ctx context.Context, name string) (models.UserEntity, error)
}

func NewUserService(log *slog.Logger,
	userStorage userStorage,
	cfg *config.Config) *UserService {
	return &UserService{
		log:         log,
		userStorage: userStorage,
		cfg:         cfg,
	}
}

func (s *UserService) GetUserByIdService(ctx context.Context, id int64) (dto.UserResponse, error) {
	op := "services.GetUserByNameService"
	log := s.log.With(slog.String("op", op))
	log.Info(op)

	userDTO := dto.UserResponse{}
	return userDTO, nil
}

func (s *UserService) GetUserByNameService(ctx context.Context, name string) (dto.UserResponse, error) {
	op := "services.GetUserByNameService"
	log := s.log.With(slog.String("op", op))
	log.Info(op)

	userDTO := dto.UserResponse{}
	return userDTO, nil
}
