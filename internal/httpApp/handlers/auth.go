package handlers

import (
	"context"
	"log/slog"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp/dto"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/gofiber/fiber/v3"
)

type authService interface {
	Register(context.Context, dto.AuthRegisterRequest) (dto.AuthRegisterResponse, error)
	Login(context.Context, dto.AuthLoginRequest) (dto.AuthLoginResponse, error)
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

func (h *AuthHandler) AuthRegister(c fiber.Ctx) error {
	op := "HttpHandlers.AuthRegister"
	log := h.log.With(slog.String("op", op))

	err := lib.ValidateBody(c, &dto.AuthRegisterRequest{})
	if err != nil {
		log.Warn(err.Error())
		return c.Status(400).SendString(err.Error())
	}

	// for i := 1; i <= 1000; i++ {
	// 	//fmt.Printf("Step %d\n", i)  // Логирование шага
	// 	//time.Sleep(1 * time.Millisecond) // Задержка 1 секунда
	// }

	body := dto.AuthRegisterRequest{}

	// 2. Используем Bind для заполнения структуры из Body
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(errorsApp.ErrInternalError.Code).JSON(fiber.Map{
			"error": "Некорректные данные: " + err.Error(),
		})
	}

	res, err := h.service.Register(c, body)
	if err != nil {
		log.Warn(err.Error())
		if err == errorsApp.ErrUserNotFound.Error {
			return c.Status(404).SendString(errorsApp.ErrUserNotFound.Message)
		}
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(200).JSON(res)
}

func (h *AuthHandler) AuthLogin(c fiber.Ctx) error {
	op := "HttpHandlers.AuthLogin"
	log := h.log.With(slog.String("op", op))

	err := lib.ValidateBody(c, &dto.AuthLoginRequest{})
	if err != nil {
		log.Warn(err.Error())
		return c.Status(400).SendString(err.Error())
	}

	// for i := 1; i <= 1000; i++ {
	// 	//fmt.Printf("Step %d\n", i)  // Логирование шага
	// 	//time.Sleep(1 * time.Millisecond) // Задержка 1 секунда
	// }

	body := dto.AuthLoginRequest{}
	//res := dto.AuthLoginResponse{}

	// 2. Используем Bind для заполнения структуры из Body
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(errorsApp.ErrInternalError.Code).JSON(fiber.Map{
			"error": "Некорректные данные: " + err.Error(),
		})
	}

	res, err := h.service.Login(c, body)
	if err != nil {
		log.Warn(err.Error())
		if err == errorsApp.ErrAuthentication.Error {
			return c.Status(401).SendString(errorsApp.ErrAuthentication.Message)
		}
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(200).JSON(res)
}
