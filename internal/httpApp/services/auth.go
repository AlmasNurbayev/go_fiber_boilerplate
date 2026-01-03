package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/config"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp/dto"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/models"
	"github.com/guregu/null/v6"
	"github.com/jinzhu/copier"
)

type AuthService struct {
	log         *slog.Logger
	authStorage authStorage
	cfg         *config.Config
}

type authStorage interface {
	NewUser(ctx context.Context, user models.UserEntity) (models.UserEntity, *errorsApp.DbError)
	GetRoleById(ctx context.Context, id int64) (models.RoleEntity, *errorsApp.DbError)
	GetUserByEmail(ctx context.Context, email string) (models.UserEntity, *errorsApp.DbError)
	GetUserByPhoneNumber(ctx context.Context, phone_number string) (models.UserEntity, *errorsApp.DbError)
	GetUserById(ctx context.Context, id int64) (models.UserEntity, *errorsApp.DbError)
}

func NewAuthService(log *slog.Logger,
	authStorage authStorage,
	cfg *config.Config) *AuthService {
	return &AuthService{
		log:         log,
		authStorage: authStorage,
		cfg:         cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, user dto.AuthRegisterRequest) (dto.AuthRegisterResponse, error) {
	op := "services.Register"
	log := s.log.With(slog.String("op", op))
	log.Info(op)

	dto := dto.AuthRegisterResponse{}

	hashedPassword, err := lib.HashPassword(user.Password)
	if err != nil {
		log.Error("error hash password", slog.String("err", err.Error()))
		return dto, err
	}

	entity, dbError := s.authStorage.NewUser(ctx, models.UserEntity{
		Name:          user.Name,
		Email:         user.Email,
		Phone_number:  user.Phone_number,
		Password_hash: null.StringFrom(hashedPassword),
		Role_id:       3, // default role user TODO - перенести в таблицу настроек
	})
	if dbError != nil {
		log.Warn("error create new user", slog.String("err", dbError.Message))
		return dto, dbError.Error
	}
	role, dbError := s.authStorage.GetRoleById(ctx, entity.Role_id)

	if dbError != nil {
		log.Warn("error create new user", slog.String("err", dbError.Message))
		return dto, dbError.Error
	}
	errCopy := copier.Copy(&dto, &entity)
	if errCopy != nil {
		log.Error("", slog.String("err", errCopy.Error()))
		return dto, errCopy
	}
	dto.Role_name = role.Name

	return dto, nil
}

func (s *AuthService) Login(ctx context.Context, user dto.AuthLoginRequest) (dto.AuthLoginResponse, error) {
	op := "services.Login"
	log := s.log.With(slog.String("op", op))
	log.Info(op)

	dto := dto.AuthLoginResponse{}
	userEntity := models.UserEntity{}

	if user.Email.Valid {
		log.Debug("login with email", slog.String("email", user.Email.String))
		userEntityByEmail, dbError := s.authStorage.GetUserByEmail(ctx, user.Email.String)
		if dbError != nil {
			if dbError.Message == "user not found" {
				log.Warn("user not found with email", slog.String("email", user.Email.String))
				return dto, errorsApp.ErrAuthentication.Error
			}

			log.Warn("error get user by email", slog.String("err", dbError.Message))
			return dto, dbError.Error
		}
		userEntity = userEntityByEmail
	}
	if user.Phone_number.Valid {
		log.Debug("login with phone number", slog.String("phone_number", user.Phone_number.String))
		userEntityByPhone, dbError := s.authStorage.GetUserByPhoneNumber(ctx, user.Phone_number.String)
		if dbError != nil {
			if dbError.Message == "user not found" {
				log.Warn("user not found with phone number", slog.String("phone_number", user.Phone_number.String))
				return dto, errorsApp.ErrAuthentication.Error
			}
			log.Warn("error get user by phone number", slog.String("err", dbError.Message))
			return dto, dbError.Error
		}
		userEntity = userEntityByPhone
	}

	err := lib.CheckPassword(userEntity.Password_hash.String, user.Password)
	if err != nil {
		log.Warn("invalid login or password", slog.String("err", err.Error()))
		return dto, errorsApp.ErrAuthentication.Error
	}

	errCopy := copier.Copy(&dto, &userEntity)
	if errCopy != nil {
		log.Error("", slog.String("err", errCopy.Error()))
		return dto, errCopy
	}
	role, dbError := s.authStorage.GetRoleById(ctx, userEntity.Role_id)

	if dbError != nil {
		log.Warn("error get role by id", slog.String("err", dbError.Message))
		return dto, dbError.Error
	}
	dto.Role_name = role.Name

	dto.AccessToken, err = lib.CreateJWT(lib.JWTClaims{
		UserId:   userEntity.Id,
		UserName: userEntity.Name,
		RoleId:   userEntity.Role_id,
		Iss:      s.cfg.SERVICE_NAME,
	}, s.cfg.AUTH_SECRET_KEY,
		time.Duration(s.cfg.AUTH_ACCESS_TOKEN_EXP_HOURS)*time.Hour,
		"access")

	if err != nil {
		log.Error("error generate access token", slog.String("err", err.Error()))
		return dto, err
	}

	dto.RefreshToken, err = lib.CreateJWT(lib.JWTClaims{
		UserId:   userEntity.Id,
		UserName: userEntity.Name,
		RoleId:   userEntity.Role_id,
		Iss:      s.cfg.SERVICE_NAME,
	}, s.cfg.AUTH_SECRET_KEY,
		time.Duration(s.cfg.AUTH_REFRESH_TOKEN_EXP_HOURS)*time.Hour,
		"refresh")

	if err != nil {
		log.Error("error generate refresh token", slog.String("err", err.Error()))
		return dto, err
	}

	return dto, nil
}

func (s *AuthService) Hello(ctx context.Context, token string) (dto.AuthHelloResponse, error) {
	op := "services.Hello"
	log := s.log.With(slog.String("op", op))
	log.Info(op)

	dto := dto.AuthHelloResponse{}
	userId, err := lib.GetUserIdFromAccessToken(token, s.cfg.AUTH_SECRET_KEY, s.cfg.SERVICE_NAME)
	if err != nil || userId == 0 {
		log.Warn("error get user id from token", slog.String("err", err.Error()))
		return dto, err
	}

	userEntity, dbError := s.authStorage.GetUserById(ctx, userId)
	if dbError != nil {
		log.Warn("error get user by id", slog.String("err", dbError.Message))
		return dto, dbError.Error
	}

	role, dbError := s.authStorage.GetRoleById(ctx, userEntity.Role_id)
	if dbError != nil {
		log.Warn("error get role by id", slog.String("err", dbError.Message))
		return dto, dbError.Error
	}

	errCopy := copier.Copy(&dto, &userEntity)
	if errCopy != nil {
		log.Error("", slog.String("err", errCopy.Error()))
		return dto, errCopy
	}
	dto.Role_name = role.Name

	return dto, nil
}

func (s *AuthService) Refresh(ctx context.Context, token string) (dto.AuthLoginResponse, error) {
	op := "services.Refresh"
	log := s.log.With(slog.String("op", op))
	log.Info(op)

	dto := dto.AuthLoginResponse{}
	claims, err := lib.GetClaimsFromRefreshToken(token, s.cfg.AUTH_SECRET_KEY, s.cfg.SERVICE_NAME)
	if err != nil || claims.UserId == 0 {
		log.Warn("error get user id from token", slog.String("err", err.Error()))
		return dto, err
	}
	dto.Id = claims.UserId
	dto.Name = claims.UserName
	dto.Role_name = "id:" + fmt.Sprint(claims.RoleId)

	dto.AccessToken, err = lib.CreateJWT(lib.JWTClaims{
		UserId:   claims.UserId,
		UserName: claims.UserName,
		RoleId:   claims.RoleId,
		Iss:      s.cfg.SERVICE_NAME,
	}, s.cfg.AUTH_SECRET_KEY,
		time.Duration(s.cfg.AUTH_ACCESS_TOKEN_EXP_HOURS)*time.Hour,
		"access")
	if err != nil {
		log.Error("internal error - generate access token", slog.String("err", err.Error()))
		return dto, err
	}

	dto.RefreshToken, err = lib.CreateJWT(lib.JWTClaims{
		UserId:   claims.UserId,
		UserName: claims.UserName,
		RoleId:   claims.RoleId,
		Iss:      s.cfg.SERVICE_NAME,
	}, s.cfg.AUTH_SECRET_KEY,
		time.Duration(s.cfg.AUTH_REFRESH_TOKEN_EXP_HOURS)*time.Hour,
		"refresh")
	if err != nil {
		log.Error("internal error - error generate refresh token", slog.String("err", err.Error()))
		return dto, err
	}

	return dto, nil
}
