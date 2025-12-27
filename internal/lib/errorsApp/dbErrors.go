package errorsApp

type DbError struct {
	Type    string
	Field   string
	Data    any
	Message string
	Error   error
}
