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
	Phone_number null.String `json:"phone_number" swaggertype:"string" example:"+77012345678"`
	Email        null.String `json:"email" swaggertype:"string" example:"almas@gmail.com"`
	Name         string      `json:"name" example:"almas"`
	Role_name    string      `json:"role_name" example:"user"`
}

type UsersResponse struct {
	Users []UserResponse `json:"users"`
}
