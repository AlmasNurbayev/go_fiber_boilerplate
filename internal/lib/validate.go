package lib

import (
	"strings"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/gofiber/fiber/v3"
)

func ValidateParams(c fiber.Ctx, dataStruct any) error {
	if err := c.Bind().URI(dataStruct); err != nil {
		return err
	}
	return nil
}

func ValidateQueryParams(c fiber.Ctx, dataStruct any) error {
	if err := c.Bind().Query(dataStruct); err != nil {
		return err
	}
	return nil
}

func ValidateBody(c fiber.Ctx, dataStruct any) error {
	if err := c.Bind().Body(dataStruct); err != nil {
		return err
	}
	return nil
}

func ExtractBearerToken(c fiber.Ctx) (string, *errorsApp.HttpError) {
	auth := c.Get("Authorization")
	if auth == "" {
		return "", &errorsApp.ErrAuthentication
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", &errorsApp.ErrAuthentication
	}

	return parts[1], nil
}
