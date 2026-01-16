package handlers

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/config"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp/dto"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/gofiber/fiber/v3"
)

type authService interface {
	Register(context.Context, dto.AuthRegisterRequest) (dto.AuthRegisterResponse, error)
	Login(context.Context, dto.AuthLoginRequest, string, string) (dto.AuthLoginResponse, error)
	Hello(context.Context, string) (dto.AuthHelloResponse, error)
	Refresh(context.Context, string) (dto.AuthLoginResponse, error)
	Sessions(context.Context, int64) (dto.AuthSessionResponse, error)
	RevokeSession(fiber.Ctx, string) error
	SendVerify(context.Context, dto.AuthSendVerifyRequest) (dto.AuthSendVerifyResponse, error)
	ConfirmVerify(context.Context, dto.AuthConfirmVerifyRequest) error
	UpdatePassword(context.Context, int64, string, string) error
}

type AuthHandler struct {
	cfg     *config.Config
	log     *slog.Logger
	service authService
}

func NewAuthHandler(cfg *config.Config, log *slog.Logger, service authService) *AuthHandler {
	return &AuthHandler{
		cfg:     cfg,
		log:     log,
		service: service,
	}
}

// @Summary      Register as user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AuthRegisterRequest  true  "Request body"
// @Success      201      {object}  dto.AuthRegisterResponse
// @Failure      409      {string}  string  "значение для поля уже существует (ограничение уникальности: users_XXXX_key)"
// @Failure      400      {string}  string  "Key: 'AuthRegisterRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag"
// @Router       /auth/register [post]
func (h *AuthHandler) AuthRegister(c fiber.Ctx) error {
	op := "HttpHandlers.AuthRegister"
	log := h.log.With(slog.String("op", op))

	err := lib.ValidateBody(c, &dto.AuthRegisterRequest{})
	if err != nil {
		log.Warn(err.Error())
		return c.Status(400).SendString(err.Error())
	}

	body := dto.AuthRegisterRequest{}

	if err := c.Bind().Body(&body); err != nil {
		return c.Status(errorsApp.ErrInternalError.Code).JSON(fiber.Map{
			"error": "Некорректные данные: " + err.Error(),
		})
	}

	res, err := h.service.Register(c, body)
	if err != nil {
		log.Warn(err.Error())
		if err == errorsApp.ErrAlreadyOtp.Error {
			return c.Status(400).SendString(errorsApp.ErrAlreadyOtp.Message)
		}
		return c.Status(400).SendString(err.Error())
	}

	return c.Status(201).JSON(res)
}

// @Summary      Login as user, returns access and refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AuthLoginRequest  true  "Request body"
// @Header       200  {string}  Set-Cookie  "refresh_token cookie is set (HttpOnly)"
// @Success      200      {object}  dto.AuthLoginResponse
// @Failure      401      {string}  string  "authentication failed"
// @Router       /auth/login [post]
func (h *AuthHandler) AuthLogin(c fiber.Ctx) error {
	op := "HttpHandlers.AuthLogin"
	log := h.log.With(slog.String("op", op))

	err := lib.ValidateBody(c, &dto.AuthLoginRequest{})
	if err != nil {
		log.Warn(err.Error())
		return c.Status(400).SendString(err.Error())
	}

	body := dto.AuthLoginRequest{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "not correct data: " + err.Error(),
		})
	}

	// sess := session.FromContext(c)
	// if err := sess.Regenerate(); err != nil { // Prevents session fixation
	// 	return err
	// }

	res, err := h.service.Login(c, body, c.IP(), string(c.Get(fiber.HeaderUserAgent)))
	if err != nil {
		log.Warn(err.Error())
		if err == errorsApp.ErrAuthentication.Error {
			return c.Status(401).SendString(errorsApp.ErrAuthentication.Message)
		}
		if err == errorsApp.ErrVerifyNotFound.Error {
			return c.Status(401).SendString(errorsApp.ErrVerifyNotFound.Message)
		}
		return c.Status(500).SendString(errorsApp.ErrInternalError.Message)
	}
	cookie := new(fiber.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = res.RefreshToken
	cookie.HTTPOnly = true
	cookie.Secure = true // true в prod (HTTPS)
	cookie.SameSite = fiber.CookieSameSiteStrictMode
	cookie.Path = "/auth/login"
	cookie.MaxAge = h.cfg.AUTH_REFRESH_TOKEN_EXP_HOURS * 60 * 60
	c.Cookie(cookie)

	return c.Status(200).JSON(res)
}

// @Summary      Check auth token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Success      200      {object}  dto.AuthHelloResponse
// @Failure      401      {string}  string  "authentication failed"
// @Router       /auth/hello [get]
func (h *AuthHandler) AuthHello(c fiber.Ctx) error {
	op := "HttpHandlers.AuthHello"
	log := h.log.With(slog.String("op", op))

	token, err := lib.ExtractBearerToken(c)
	if err != nil {
		log.Warn(err.Message)
		return c.Status(err.Code).SendString(err.Message)
	}

	res, err2 := h.service.Hello(c, token)
	if err2 != nil {
		log.Warn(err2.Error())
		return c.Status(errorsApp.ErrAuthentication.Code).SendString(errorsApp.ErrAuthentication.Message)
	}
	return c.Status(200).JSON(res)
}

// @Summary      Check refresh token, returns new access and refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Success      200      {object}  dto.AuthLoginResponse
// @Header       200  {string}  Set-Cookie  "refresh_token cookie is set (HttpOnly)"
// @Failure      401      {string}  string  "authentication failed"
// @Router       /auth/refresh [post]
func (h *AuthHandler) AuthRefresh(c fiber.Ctx) error {
	op := "HttpHandlers.AuthRefresh"
	log := h.log.With(slog.String("op", op))

	token, err := lib.ExtractBearerToken(c)
	if err != nil {
		log.Warn(err.Message)
		return c.Status(err.Code).SendString(err.Message)
	}

	res, err2 := h.service.Refresh(c, token)
	if err2 != nil {
		log.Warn(err2.Error())
		if strings.Contains(err2.Error(), "internal error") {
			return c.Status(500).SendString(errorsApp.ErrInternalError.Message)
		}
		if err2 == errorsApp.ErrSessionNotFound.Error {
			return c.Status(401).SendString(errorsApp.ErrSessionNotFound.Message)
		}

		return c.Status(401).SendString(errorsApp.ErrAuthentication.Message)
	}
	cookie := new(fiber.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = res.RefreshToken
	cookie.HTTPOnly = true
	cookie.SameSite = fiber.CookieSameSiteStrictMode
	cookie.Path = "/auth/refresh"
	cookie.MaxAge = h.cfg.AUTH_REFRESH_TOKEN_EXP_HOURS * 60 * 60
	c.Cookie(cookie)

	return c.Status(200).JSON(res)
}

// @Summary      Get sessions by user id, only user can get his own sessions
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Success      200      {object}  dto.AuthSessionResponse
// @Param        id  path      string  true  "User id"
// @Failure      401      {string}  string  "authentication failed"
// @Router       /auth/sessions/{id} [get]
func (h *AuthHandler) AuthSessions(c fiber.Ctx) error {
	op := "HttpHandlers.AuthSessions"
	log := h.log.With(slog.String("op", op))

	idString := c.Params("id")
	userId, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		log.Warn(err.Error())
		return c.Status(400).SendString(err.Error())
	}

	// проверяем этот же ли запршивает
	userIdFromContext := c.Locals("user_id").(int64)
	if userId != userIdFromContext {
		log.Warn("user id not match", slog.String("err", "user id not match"))
		return c.Status(errorsApp.ErrForbidden.Code).SendString(errorsApp.ErrForbidden.Message)
	}

	res := dto.AuthSessionResponse{}

	res, err2 := h.service.Sessions(c, userId)
	if err2 != nil {
		log.Warn(err2.Error())
		if strings.Contains(err2.Error(), "internal error") {
			return c.Status(500).SendString(errorsApp.ErrInternalError.Message)
		}
		if err2 == errorsApp.ErrSessionNotFound.Error {
			return c.Status(401).SendString(errorsApp.ErrSessionNotFound.Message)
		}

		return c.Status(401).SendString(errorsApp.ErrAuthentication.Message)
	}

	return c.Status(200).JSON(res)
}

// @Summary      Revoke session by jti, only user can revoke his own session
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Success      200      string  "ok"
// @Param        jti  path      string  true  "Session jti"
// @Failure      401      {string}  string  "authentication failed"
// @Router       /auth/sessions/{jti} [delete]
func (h *AuthHandler) RevokeSession(c fiber.Ctx) error {
	op := "HttpHandlers.RevokeSession"
	log := h.log.With(slog.String("op", op))

	jtiString := c.Params("jti")

	err := h.service.RevokeSession(c, jtiString)
	if err != nil {
		log.Warn(err.Error())
		if strings.Contains(err.Error(), "internal_error") {
			return c.Status(500).SendString(errorsApp.ErrInternalError.Message)
		}
		if err == errorsApp.ErrSessionNotFound.Error {
			return c.Status(401).SendString(errorsApp.ErrSessionNotFound.Message)
		}
		if err == errorsApp.ErrForbidden.Error {
			return c.Status(403).SendString(errorsApp.ErrForbidden.Message)
		}

		return c.Status(401).SendString(errorsApp.ErrAuthentication.Message)
	}

	return c.Status(200).SendString("ok")
}
