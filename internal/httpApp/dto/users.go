package dto

import "github.com/guregu/null/v6"

type UserRequestParams struct {
	Id int64 `validate:"required,gte=0" example:"5"`
}

type UserRequestQueryParams struct {
	Name string `query:"name" validate:"omitempty" example:"almas"`
}

type UserResponse struct {
	Id           int64       `json:"id"`
	Phone_number null.String `json:"phone_number"`
	Email        null.String `json:"email"`
	Name         string      `json:"name"`
	Role_name    string      `json:"role_name"`
}
