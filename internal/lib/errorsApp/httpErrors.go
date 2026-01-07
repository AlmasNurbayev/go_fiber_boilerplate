package errorsApp

import "errors"

type HttpError struct {
	Code    int
	Message string
	Error   error
}

var (
	ErrTimeout = HttpError{
		Code:    408,
		Message: "time out",
		Error:   errors.New("time out")}

	ErrUserNotFound = HttpError{
		Code:    404,
		Message: "user not found",
		Error:   errors.New("user not found")}

	ErrInternalError = HttpError{
		Code:    500,
		Message: "internal error",
		Error:   errors.New("internal error")}

	ErrBadRequest = HttpError{
		Code:    400,
		Message: "bad request",
		Error:   errors.New("bad request")}

	ErrNewsNotFound = HttpError{
		Code:    404,
		Message: "news not found",
		Error:   errors.New("news not found")}

	ErrMaxPriceLessMinPrice = HttpError{
		Code:    400,
		Message: "max price less then min price",
		Error:   errors.New("max price less then min price")}

	ErrSortBadFormat = HttpError{
		Code:    400,
		Message: "sort don't contain -",
		Error:   errors.New("sort don't contain -")}

	ErrProductNotFound = HttpError{
		Code:    404,
		Message: "product not found",
		Error:   errors.New("product not found")}

	ErrKaspiCategoryDuplicate = HttpError{
		Code:    400,
		Message: "kaspi category is exists",
		Error:   errors.New("kaspi category is exists")}

	ErrAuthentication = HttpError{
		Code:    401,
		Message: "authentication failed",
		Error:   errors.New("authentication failed")}

	ErrSessionNotFound = HttpError{
		Code:    401,
		Message: "session not found or expired",
		Error:   errors.New("session not found or expired")}
)
