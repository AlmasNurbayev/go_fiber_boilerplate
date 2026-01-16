package handlers

import (
	"log/slog"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp/dto"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib"
	"github.com/gofiber/fiber/v3"
)

// @Summary      Verify user addresses
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AuthSendVerifyRequest  true  "Request body"
// @Success      200      {object}  dto.AuthSendVerifyResponse
// @Failure      401      {string}  string  "authentication failed"
// @Router       /auth/send-verify [post]
func (h *AuthHandler) SendVerify(c fiber.Ctx) error {
	op := "HttpHandlers.SendVerify"
	log := h.log.With(slog.String("op", op))

	err := lib.ValidateBody(c, &dto.AuthSendVerifyRequest{})
	if err != nil {
		log.Warn(err.Error())
		return c.Status(400).SendString(err.Error())
	}

	body := dto.AuthSendVerifyRequest{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "not correct data: " + err.Error(),
		})
	}

	responseOtp, err2 := h.service.SendVerify(c, body)
	if err2 != nil {
		log.Warn(err2.Error())
		return c.Status(400).SendString(err2.Error())
	}

	return c.Status(200).JSON(responseOtp)
}

// @Summary      Send verify code to user address
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AuthConfirmVerifyRequest  true  "Request body"
// @Success      200      string  "ok"
// @Failure      400      {string}  string  "bad request"
// @Router       /auth/confirm-verify [post]
func (h *AuthHandler) ConfirmVerify(c fiber.Ctx) error {
	op := "HttpHandlers.ConfirmVerify"
	log := h.log.With(slog.String("op", op))

	err := lib.ValidateBody(c, &dto.AuthConfirmVerifyRequest{})
	if err != nil {
		log.Warn(err.Error())
		return c.Status(400).SendString(err.Error())
	}

	body := dto.AuthConfirmVerifyRequest{}
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "not correct data: " + err.Error(),
		})
	}

	err2 := h.service.ConfirmVerify(c, body)
	if err2 != nil {
		log.Warn(err2.Error())
		return c.Status(400).SendString(err2.Error())
	}

	return c.Status(200).SendString("ok")
}
