package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/config"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/db/cache"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp/dto"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/guregu/null/v6"
	"github.com/jinzhu/copier"
)

type AuthService struct {
	log            *slog.Logger
	authStorage    authStorage
	sessionStorage sessionStorage
	otpStorage     otpStorage
	cfg            *config.Config
}

type authStorage interface {
	NewUser(ctx context.Context, user models.UserEntity) (models.UserEntity, *errorsApp.DbError)
	GetRoleById(ctx context.Context, id int64) (models.RoleEntity, *errorsApp.DbError)
	GetUserByEmail(ctx context.Context, email string) (models.UserEntity, *errorsApp.DbError)
	GetUserByPhoneNumber(ctx context.Context, phone_number string) (models.UserEntity, *errorsApp.DbError)
	GetUserById(ctx context.Context, id int64) (models.UserEntity, *errorsApp.DbError)
	UpdateUserEmailVerifyTimestamp(ctx context.Context, id int64) *errorsApp.DbError
	UpdateUserPhoneVerifyTimestamp(ctx context.Context, id int64) *errorsApp.DbError
	UpdatePassword(ctx context.Context, id int64, password string) *errorsApp.DbError
}

type sessionStorage interface {
	SaveSession(ctx context.Context, jti string, data cache.SessionData, ttlHours int) *errorsApp.DbError
	GetSessionByJti(ctx context.Context, jti string) (cache.SessionData, *errorsApp.DbError)
	GetSessionsByUserId(ctx context.Context, userId int64) ([]cache.SessionData, *errorsApp.DbError)
	DeleteSessionByJti(ctx context.Context, jti string) *errorsApp.DbError
}

type otpStorage interface {
	SaveOtp(ctx context.Context, data cache.OtpData, ttlMinutes int) *errorsApp.DbError
	DeleteOtp(ctx context.Context, address string, typeM string) *errorsApp.DbError
	GetOtp(ctx context.Context, address string, typeM string) (cache.OtpData, *errorsApp.DbError)
}

func NewAuthService(log *slog.Logger,
	authStorage authStorage,
	sessionStorage sessionStorage,
	otpStorage otpStorage,
	cfg *config.Config) *AuthService {
	return &AuthService{
		log:            log,
		authStorage:    authStorage,
		sessionStorage: sessionStorage,
		otpStorage:     otpStorage,
		cfg:            cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, user dto.AuthRegisterRequest) (dto.AuthRegisterResponse, error) {
	op := "services.Register"
	log := s.log.With(slog.String("op", op))

	response := dto.AuthRegisterResponse{}
	hashedPassword, err := lib.HashPassword(user.Password)
	if err != nil {
		log.Error("error hash password", slog.String("err", err.Error()))
		return response, err
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
		return response, dbError.Error
	}
	role, dbError := s.authStorage.GetRoleById(ctx, entity.Role_id)

	if dbError != nil {
		log.Warn("error create new user", slog.String("err", dbError.Message))
		return response, dbError.Error
	}
	errCopy := copier.Copy(&response, &entity)
	if errCopy != nil {
		log.Error("", slog.String("err", errCopy.Error()))
		return response, errCopy
	}
	response.Role_name = role.Name

	// отправляем код подтверждения
	switch user.ConfirmType {
	case "phone":
		responseSend, errSendVerify := s.SendVerify(ctx, dto.AuthSendVerifyRequest{
			Type:    "phone",
			Address: user.Phone_number.String,
		})
		if errSendVerify != nil {
			log.Warn("error send verify", slog.String("err", errSendVerify.Error()))
			return response, errSendVerify
		}
		response.OtpExpiresAt = responseSend.OtpExpiresAt
	case "email":
		responseSend, errSendVerify := s.SendVerify(ctx, dto.AuthSendVerifyRequest{
			Type:    "email",
			Address: user.Email.String,
		})
		if errSendVerify != nil {
			if errSendVerify == errorsApp.ErrAlreadyOtp.Error {
				return response, errorsApp.ErrAlreadyOtp.Error
			}
			log.Warn("error send verify", slog.String("err", errSendVerify.Error()))
			return response, errSendVerify
		}
		response.OtpExpiresAt = responseSend.OtpExpiresAt
	default:
		log.Warn("invalid confirm type", slog.String("confirm_type", user.ConfirmType))
		return response, errorsApp.ErrBadRequest.Error
	}
	return response, nil
}

func (s *AuthService) Login(ctx context.Context, user dto.AuthLoginRequest, ip string, user_agent string) (dto.AuthLoginResponse, error) {
	op := "services.Login"
	log := s.log.With(slog.String("op", op))

	dto := dto.AuthLoginResponse{}
	userEntity := models.UserEntity{}

	// проверяем наличие пользователя по email
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
		// проверяем наличие верификацию пользователя по email
		if !userEntityByEmail.Email_verified_at.Valid {
			log.Warn("user not verified", slog.String("name", userEntityByEmail.Name))
			return dto, errorsApp.ErrVerifyNotFound.Error
		}
		userEntity = userEntityByEmail
	}

	// проверяем наличие пользователя по номеру телефона
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
		// проверяем наличие верификацию пользователя по email
		if !userEntityByPhone.Email_verified_at.Valid {
			log.Warn("user not verified", slog.String("name", userEntityByPhone.Name))
			return dto, errorsApp.ErrVerifyNotFound.Error
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

	jti := uuid.New().String()

	dto.AccessToken, err = lib.CreateJWT(lib.JWTClaims{
		UserId:   userEntity.Id,
		UserName: userEntity.Name,
		RoleId:   userEntity.Role_id,
		Jti:      jti,
		Iss:      s.cfg.SERVICE_NAME,
	}, s.cfg.AUTH_SECRET_KEY,
		time.Duration(s.cfg.AUTH_ACCESS_TOKEN_EXP_MINUTES)*time.Minute,
		"access")

	if err != nil {
		log.Error("error generate access token", slog.String("err", err.Error()))
		return dto, err
	}

	dto.RefreshToken, err = lib.CreateJWT(lib.JWTClaims{
		UserId:   userEntity.Id,
		UserName: userEntity.Name,
		RoleId:   userEntity.Role_id,
		Jti:      jti,
		Iss:      s.cfg.SERVICE_NAME,
	}, s.cfg.AUTH_SECRET_KEY,
		time.Duration(s.cfg.AUTH_REFRESH_TOKEN_EXP_HOURS)*time.Hour,
		"refresh")

	if err != nil {
		log.Error("error generate refresh token", slog.String("err", err.Error()))
		return dto, err
	}

	err2 := s.sessionStorage.SaveSession(ctx, jti, cache.SessionData{
		Jti:       jti,
		UserID:    userEntity.Id,
		RoleID:    userEntity.Role_id,
		UserAgent: user_agent,
		IP:        ip,
	}, s.cfg.AUTH_REFRESH_TOKEN_EXP_HOURS)
	if err2 != nil {
		log.Error("error save session", slog.String("err", err2.Message))
		return dto, err2.Error
	}

	return dto, nil
}

func (s *AuthService) Hello(ctx context.Context, token string) (dto.AuthHelloResponse, error) {
	op := "services.Hello"
	log := s.log.With(slog.String("op", op))

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

	dto := dto.AuthLoginResponse{}
	claims, err := lib.GetClaimsFromRefreshToken(token, s.cfg.AUTH_SECRET_KEY, s.cfg.SERVICE_NAME)
	if err != nil || claims.UserId == 0 {
		log.Warn("error get user id from token", slog.String("err", err.Error()))
		return dto, err
	}

	//Проверяем наличие JTI в Redis (Whitelist)
	data, err2 := s.sessionStorage.GetSessionByJti(ctx, claims.Jti)
	if err2 != nil {
		log.Warn("error get session by jti", slog.String("err", err2.Message))
		return dto, errorsApp.ErrSessionNotFound.Error
	}
	if data.UserID != claims.UserId {
		log.Warn("refresh-token user_id not match session user_id", slog.Int64("user_id", data.UserID), slog.Int64("claims_user_id", claims.UserId))
		return dto, errorsApp.ErrSessionNotFound.Error
	}

	// удаляем старую сессию
	err3 := s.sessionStorage.DeleteSessionByJti(ctx, "jti:"+claims.Jti)
	if err3 != nil {
		log.Warn("error delete session by jti", slog.String("err", err3.Message))
		// ничего не делаем если не удалось удалить сессию
	}

	// сохраняем новую сессию с другим jti
	newJti := uuid.NewString()
	err4 := s.sessionStorage.SaveSession(ctx, newJti, data, s.cfg.AUTH_REFRESH_TOKEN_EXP_HOURS)
	if err4 != nil {
		log.Error("error save session by jti", slog.String("err", err4.Message))
		return dto, errorsApp.ErrInternalError.Error
	}

	dto.Id = claims.UserId
	dto.Name = claims.UserName
	dto.Role_name = "id:" + fmt.Sprint(claims.RoleId)

	dto.AccessToken, err = lib.CreateJWT(lib.JWTClaims{
		UserId:   claims.UserId,
		UserName: claims.UserName,
		Jti:      newJti,
		RoleId:   claims.RoleId,
		Iss:      s.cfg.SERVICE_NAME,
	}, s.cfg.AUTH_SECRET_KEY,
		time.Duration(s.cfg.AUTH_ACCESS_TOKEN_EXP_MINUTES)*time.Minute,
		"access")
	if err != nil {
		log.Error("internal error - generate access token", slog.String("err", err.Error()))
		return dto, err
	}

	dto.RefreshToken, err = lib.CreateJWT(lib.JWTClaims{
		UserId:   claims.UserId,
		UserName: claims.UserName,
		RoleId:   claims.RoleId,
		Jti:      newJti,
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

func (s *AuthService) Sessions(ctx context.Context, id int64) (dto.AuthSessionResponse, error) {
	op := "services.Sessions"
	log := s.log.With(slog.String("op", op))

	response := dto.AuthSessionResponse{}
	response.Sessions = make([]dto.AuthSession, 0)

	sessionData, err2 := s.sessionStorage.GetSessionsByUserId(ctx, id)
	if err2 != nil {
		log.Warn("error get sessions by user id", slog.String("err", err2.Message))
		return response, errorsApp.ErrSessionNotFound.Error
	}

	userData, err3 := s.authStorage.GetUserById(ctx, id)
	if err3 != nil {
		log.Warn("error get user by id", slog.Any("err", err3))
		return response, errorsApp.ErrUserNotFound.Error
	}

	for _, session := range sessionData {
		response.Sessions = append(response.Sessions, dto.AuthSession{
			Jti:               session.Jti,
			User_id:           session.UserID,
			User_name:         userData.Name,
			User_email:        userData.Email.String,
			User_phone_number: userData.Phone_number.String,
			Role_id:           session.RoleID,
			User_agent:        session.UserAgent,
			IP:                session.IP,
			Created_at:        session.CreatedAt,
		})
	}

	return response, nil
}

func (s *AuthService) RevokeSession(ctx fiber.Ctx, jtiString string) error {
	op := "services.RevokeSession"
	log := s.log.With(slog.String("op", op))

	userId := ctx.Locals("user_id").(int64)
	if userId == 0 {
		log.Warn("user id not found", slog.String("err", "user id not found"))
		return errorsApp.ErrAuthentication.Error
	}
	data, err := s.sessionStorage.GetSessionByJti(ctx, jtiString)
	if err != nil {
		log.Warn("error get session by jti", slog.String("err", err.Message))
		return errorsApp.ErrAuthentication.Error
	}
	if data.UserID != userId {
		log.Warn("jti user_id not match session user_id", slog.String("err", "jti user_id not match session user_id"))
		return errorsApp.ErrForbidden.Error
	}

	err2 := s.sessionStorage.DeleteSessionByJti(ctx, "jti:"+jtiString)
	if err2 != nil {
		log.Warn("error delete session by jti", slog.String("err", err2.Message))
		switch err2.Type {
		case "not_found":
			return errorsApp.ErrSessionNotFound.Error
		case "internal_error":
			return errorsApp.ErrInternalError.Error
		default:
			return errorsApp.ErrAuthentication.Error
		}
	}

	return nil
}

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

	return nil
}
