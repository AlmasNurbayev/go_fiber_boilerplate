package handlers

import (
	"context"
	"log/slog"
	"strings"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp/dto"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/gofiber/fiber/v3"
)

type authService interface {
	Register(context.Context, dto.AuthRegisterRequest) (dto.AuthRegisterResponse, error)
	Login(context.Context, dto.AuthLoginRequest) (dto.AuthLoginResponse, error)
	Hello(context.Context, string) (dto.AuthHelloResponse, error)
	Refresh(context.Context, string) (dto.AuthLoginResponse, error)
}

type AuthHandler struct {
	log     *slog.Logger
	service authService
}

func NewAuthHandler(log *slog.Logger, service authService) *AuthHandler {
	return &AuthHandler{
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
		return c.Status(400).SendString(err.Error())
	}

	return c.Status(201).JSON(res)
}

// @Summary      Login as user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AuthLoginRequest  true  "Request body"
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

	res, err := h.service.Login(c, body)
	if err != nil {
		log.Warn(err.Error())
		if err == errorsApp.ErrAuthentication.Error {
			return c.Status(401).SendString(errorsApp.ErrAuthentication.Message)
		}
		return c.Status(500).SendString(errorsApp.ErrInternalError.Message)
	}
	//sess.Set("authenticated", true)

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

// @Summary      Refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Success      200      {object}  dto.AuthLoginResponse
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
		return c.Status(401).SendString(errorsApp.ErrAuthentication.Message)
	}

	return c.Status(200).JSON(res)
}
