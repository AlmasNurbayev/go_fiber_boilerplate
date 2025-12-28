package dto

import "github.com/guregu/null/v6"

type AuthRegisterRequest struct {
	Phone_number null.String `json:"phone_number" validate:"required_without=Email,omitempty" swaggertype:"string" example:"+77012345678"`
	Email        null.String `json:"email" validate:"required_without=Phone_number,omitempty" swaggertype:"string" example:"test@mail.com"`
	Name         string      `json:"name" validate:"required"`
	Password     string      `json:"password" validate:"required,min=8"`
}

type AuthRegisterResponse struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Role_name string `json:"role_name"`
}

type AuthLoginRequest struct {
	Phone_number null.String `json:"phone_number" validate:"required_without=Email,omitempty" swaggertype:"string" example:"+77012345678"`
	Email        null.String `json:"email" validate:"required_without=Phone_number,omitempty" swaggertype:"string" example:"test@mail.com"`
	Password     string      `json:"password" validate:"required,min=8"`
}

type AuthLoginResponse struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Role_name   string `json:"role_name"`
	AccessToken string `json:"access_token"`
}

type AuthHelloResponse struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Role_name    string `json:"role_name"`
	Email        string `json:"email"`
	Phone_number string `json:"phone_number"`
}
