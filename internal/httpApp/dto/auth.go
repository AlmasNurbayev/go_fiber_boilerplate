package dto

import (
	"time"

	"github.com/guregu/null/v6"
)

type AuthRegisterRequest struct {
	Phone_number null.String `json:"phone_number" validate:"required_without=Email,omitempty,phoneKZ" swaggertype:"string" example:"77012345678"`
	Email        null.String `json:"email" validate:"required_without=Phone_number,omitempty" swaggertype:"string" example:"test@mail.com"`
	Name         string      `json:"name" validate:"required"`
	ConfirmType  string      `json:"confirm_type" validate:"required" swaggertype:"string" example:"phone or email"`
	Password     string      `json:"password" validate:"required,min=8"`
}

type AuthRegisterResponse struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	Role_name    string    `json:"role_name"`
	OtpExpiresAt time.Time `json:"otp_expires_at"`
}

type AuthLoginRequest struct {
	Phone_number null.String `json:"phone_number" validate:"required_without=Email,omitempty,phoneKZ" swaggertype:"string" example:"+77012345678"`
	Email        null.String `json:"email" validate:"required_without=Phone_number,omitempty" swaggertype:"string" example:"test@mail.com"`
	Password     string      `json:"password" validate:"required,min=8"`
}

type AuthLoginResponse struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Role_name    string `json:"role_name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthHelloResponse struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Role_name    string `json:"role_name"`
	Email        string `json:"email"`
	Phone_number string `json:"phone_number"`
}

type AuthSession struct {
	Jti               string    `json:"jti"`
	User_id           int64     `json:"user_id"`
	User_name         string    `json:"user_name"`
	User_email        string    `json:"user_email"`
	User_phone_number string    `json:"user_phone_number"`
	Role_id           int64     `json:"role_id"`
	User_agent        string    `json:"user_agent"`
	IP                string    `json:"ip"`
	Created_at        time.Time `json:"created_at"`
}

type AuthSessionResponse struct {
	Sessions []AuthSession `json:"sessions"`
}

type AuthSendVerifyRequest struct {
	Type    string `json:"type" validate:"required" swaggertype:"string" example:"phone"`
	Address string `json:"address" validate:"required" swaggertype:"string" example:"+77012345678"`
}

type AuthConfirmVerifyRequest struct {
	//UserID  int64  `json:"user_id" validate:"required" swaggertype:"integer" example:"1"`
	Type    string `json:"type" validate:"required" swaggertype:"string" example:"phone"`
	Address string `json:"address" validate:"required" swaggertype:"string" example:"+77012345678"`
	Code    string `json:"code" validate:"required,min=6,max=6"`
}

type AuthSendVerifyResponse struct {
	OtpExpiresAt time.Time `json:"otp_expires_at"`
}

type AuthUpdatePasswordRequest struct {
	UserId      int64  `json:"user_id" validate:"required" swaggertype:"integer" example:"1"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
