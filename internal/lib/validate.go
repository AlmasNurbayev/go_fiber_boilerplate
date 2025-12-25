package lib

import (
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
