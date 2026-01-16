package handlers

import (
	"log/slog"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp/dto"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/gofiber/fiber/v3"
)

// @Summary      Update user password for authenticated user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.AuthUpdatePasswordRequest  true  "Request body"
// @Success      200      string  "ok"
// @Failure      400      {string}  string  "bad request"
// @Failure      401      {string}  string  "authentication failed"
// @Router       /auth/update-password [post]
func (h *AuthHandler) UpdatePassword(c fiber.Ctx) error {
	op := "HttpHandlers.UpdatePassword"
	log := h.log.With(slog.String("op", op))

	err := lib.ValidateBody(c, &dto.AuthUpdatePasswordRequest{})
	if err != nil {
		log.Warn(err.Error())
		return c.Status(400).SendString(err.Error())
	}

	body := dto.AuthUpdatePasswordRequest{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "not correct data: " + err.Error(),
		})
	}

	err2 := h.service.UpdatePassword(c, body.UserId, body.OldPassword, body.NewPassword)
	if err2 != nil {
		if err2 == errorsApp.ErrAuthentication.Error {
			return c.Status(401).SendString(err2.Error())
		}
		if err2 == errorsApp.ErrOldPasswordNotMatch.Error {
			return c.Status(errorsApp.ErrOldPasswordNotMatch.Code).SendString(err2.Error())
		}
		log.Warn(err2.Error())
		return c.Status(400).SendString(err2.Error())
	}

	return c.Status(200).SendString("ok")
}
