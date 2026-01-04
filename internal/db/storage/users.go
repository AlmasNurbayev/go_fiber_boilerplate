package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (s *Storage) GetUserByNameStorage(ctx context.Context, name string) ([]models.UserEntity, error) {
	op := "storage.GetUserByNameStorage"
	log := s.log.With("op", op)

	var users = []models.UserEntity{}

	query := `SELECT id, name, email, role_id FROM "users" WHERE name = $1`
	err := pgxscan.Select(ctx, s.Db, &users, query, name)
	if err != nil {
		log.Error(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return users, errorsApp.ErrUserNotFound.Error
		}
		return users, errorsApp.ErrInternalError.Error
	}

	return users, nil
}
