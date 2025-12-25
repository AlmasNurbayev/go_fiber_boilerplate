package lib

import "errors"

type ErrorHttp struct {
	Code    int
	Message string
	Error   error
}

var (
	ErrTimeout = ErrorHttp{
		Code:    408,
		Message: "time out",
		Error:   errors.New("time out")}

	ErrUserNotFound = ErrorHttp{
		Code:    404,
		Message: "user not found",
		Error:   errors.New("user not found")}

	ErrInternalError = ErrorHttp{
		Code:    500,
		Message: "internal error",
		Error:   errors.New("internal error")}

	ErrBadRequest = ErrorHttp{
		Code:    400,
		Message: "bad request",
		Error:   errors.New("bad request")}

	ErrNewsNotFound = ErrorHttp{
		Code:    404,
		Message: "news not found",
		Error:   errors.New("news not found")}

	ErrMaxPriceLessMinPrice = ErrorHttp{
		Code:    400,
		Message: "max price less then min price",
		Error:   errors.New("max price less then min price")}

	ErrSortBadFormat = ErrorHttp{
		Code:    400,
		Message: "sort don't contain -",
		Error:   errors.New("sort don't contain -")}

	ErrProductNotFound = ErrorHttp{
		Code:    404,
		Message: "product not found",
		Error:   errors.New("product not found")}

	ErrKaspiCategoryDuplicate = ErrorHttp{
		Code:    400,
		Message: "kaspi category is exists",
		Error:   errors.New("kaspi category is exists")}
)
