package errorsApp

import "errors"

type httpError struct {
	Code    int
	Message string
	Error   error
}

var (
	ErrTimeout = httpError{
		Code:    408,
		Message: "time out",
		Error:   errors.New("time out")}

	ErrUserNotFound = httpError{
		Code:    404,
		Message: "user not found",
		Error:   errors.New("user not found")}

	ErrInternalError = httpError{
		Code:    500,
		Message: "internal error",
		Error:   errors.New("internal error")}

	ErrBadRequest = httpError{
		Code:    400,
		Message: "bad request",
		Error:   errors.New("bad request")}

	ErrNewsNotFound = httpError{
		Code:    404,
		Message: "news not found",
		Error:   errors.New("news not found")}

	ErrMaxPriceLessMinPrice = httpError{
		Code:    400,
		Message: "max price less then min price",
		Error:   errors.New("max price less then min price")}

	ErrSortBadFormat = httpError{
		Code:    400,
		Message: "sort don't contain -",
		Error:   errors.New("sort don't contain -")}

	ErrProductNotFound = httpError{
		Code:    404,
		Message: "product not found",
		Error:   errors.New("product not found")}

	ErrKaspiCategoryDuplicate = httpError{
		Code:    400,
		Message: "kaspi category is exists",
		Error:   errors.New("kaspi category is exists")}

	ErrAuthentication = httpError{
		Code:    401,
		Message: "authentication failed",
		Error:   errors.New("authentication failed")}
)
