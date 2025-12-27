package storage

import (
	"errors"
	"fmt"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/jackc/pgx/v5/pgconn"
)

func mapPgError(err error) *errorsApp.DbError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			// Можно проверить pgErr.ConstraintName, чтобы понять, какое именно поле дублируется
			return &errorsApp.DbError{
				Type:    "unique_violation",
				Field:   pgErr.ConstraintName,
				Message: fmt.Sprintf("значение для поля уже существует (ограничение уникальности: %s)", pgErr.ConstraintName),
				Error:   fmt.Errorf("значение для поля уже существует (ограничение уникальности: %s)", pgErr.ConstraintName),
			}

		case "23502": // not_null_violation
			return &errorsApp.DbError{
				Type:    "not_null_violation",
				Field:   pgErr.ColumnName,
				Message: fmt.Sprintf("пропущено обязательное поле: %s", pgErr.ColumnName),
				Error:   fmt.Errorf("пропущено обязательное поле: %s", pgErr.ColumnName),
			}

		case "23503": // foreign_key_violation
			return &errorsApp.DbError{
				Type:    "foreign_key_violation",
				Field:   pgErr.ConstraintName,
				Message: fmt.Sprintf("ошибка внешнего ключа: %s (проверьте ограничение %s)", pgErr.Detail, pgErr.ConstraintName),
				Error:   fmt.Errorf("ошибка внешнего ключа: %s (проверьте ограничение %s)", pgErr.Detail, pgErr.ConstraintName),
			}

		}
	}
	return &errorsApp.DbError{
		Type:    "unknown_database_error",
		Message: "внутренняя ошибка базы данных",
		Error:   fmt.Errorf("внутренняя ошибка базы данных: %s", pgErr.Message),
	}
}
