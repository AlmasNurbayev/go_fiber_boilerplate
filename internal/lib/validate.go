package lib

import (
	"reflect"
	"strings"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/guregu/null/v6"
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

func isValidPhoneKZ(phone string) bool {
	if len(phone) != 11 {
		return false
	}
	if phone[0] != '7' {
		return false
	}
	for i := range 11 {
		if phone[i] < '0' || phone[i] > '9' {
			return false
		}
	}
	return true
}

func PhoneValidatorKZ(fl validator.FieldLevel) bool {
	switch fl.Field().Type() {
	case reflect.TypeFor[string]():
		phone := fl.Field().String()
		return isValidPhoneKZ(phone)
	case reflect.TypeFor[null.String]():
		ns := fl.Field().Interface().(null.String)
		if !ns.Valid {
			return false
		}
		return isValidPhoneKZ(ns.String)
	default:
		return false
	}
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
