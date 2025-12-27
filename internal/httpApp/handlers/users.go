package handlers

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp/dto"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/gofiber/fiber/v3"
)

type userServices interface {
	GetUserByIdService(ctx context.Context, id int64) (dto.UserResponse, error)
	GetUserByNameService(ctx context.Context, name string) (dto.UserResponse, error)
}

type UserHandler struct {
	log     *slog.Logger
	service userServices
}

func NewUserHandler(log *slog.Logger, service userServices) *UserHandler {
	return &UserHandler{
		log:     log,
		service: service,
	}
}

func (h *UserHandler) GetUserById(c fiber.Ctx) error {
	op := "HttpHandlers.GetUser"
	log := h.log.With(slog.String("op", op))

	err := lib.ValidateParams(c, &dto.UserRequestParams{})
	if err != nil {
		log.Warn(err.Error())
		return c.Status(400).SendString(err.Error())
	}

	for i := 1; i <= 2000; i++ {
		//fmt.Printf("Step %d\n", i)  // Логирование шага
		//time.Sleep(1 * time.Millisecond) // Задержка 1 секунда
	}

	idString := c.Params("id")
	res := dto.UserResponse{}

	if idString != "" {
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			log.Warn(err.Error())
			return c.Status(400).SendString(errorsApp.ErrBadRequest.Message)
		}

		res, err = h.service.GetUserByIdService(c, id)
		if err != nil {
			log.Warn(err.Error())
			if err == errorsApp.ErrUserNotFound.Error {
				return c.Status(404).SendString(errorsApp.ErrUserNotFound.Message)
			}
			return c.Status(500).SendString(err.Error())
		}

	}
	return c.Status(200).JSON(res)
}

func (h *UserHandler) GetUserSearch(c fiber.Ctx) error {
	op := "HttpHandlers.GetUserSearch"
	log := h.log.With(slog.String("op", op))

	err := lib.ValidateQueryParams(c, &dto.UserRequestQueryParams{})
	if err != nil {
		log.Warn(err.Error())
		return c.Status(400).SendString(err.Error())
	}

	for i := 1; i <= 2000; i++ {
		//fmt.Printf("Step %d\n", i)  // Логирование шага
		//time.Sleep(1 * time.Millisecond) // Задержка 1 секунда
	}

	nameString := c.Query("name")
	res := dto.UserResponse{}

	if nameString != "" {
		res, err = h.service.GetUserByNameService(c, nameString)
		if err != nil {
			log.Warn(err.Error())
			if err == errorsApp.ErrUserNotFound.Error {
				return c.Status(404).SendString(errorsApp.ErrUserNotFound.Message)
			}
			return c.Status(500).SendString(err.Error())
		}
	}
	return c.Status(200).JSON(res)
}
