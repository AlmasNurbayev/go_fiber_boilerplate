package dto

import (
	"github.com/guregu/null/v6"
)

type UserRequestParams struct {
	Id int64 `validate:"required,gte=0" example:"5"`
}

type UserRequestQueryParams struct {
	Name string `query:"name" validate:"omitempty" example:"almas"`
}

type UserResponse struct {
	Id           int64       `json:"id"`
	Phone_number string      `json:"phone_number"`
	Email        string      `json:"email"`
	Name         null.String `json:"name"`
	Role_name    null.String `json:"role_name"`
}
