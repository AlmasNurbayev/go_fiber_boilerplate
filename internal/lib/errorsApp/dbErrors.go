package errorsApp

type DbError struct {
	Type    string // not_found, internal_error
	Field   string
	Data    any
	Message string
	Error   error
}
