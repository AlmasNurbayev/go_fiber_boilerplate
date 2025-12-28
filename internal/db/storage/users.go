package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (s *Storage) GetUserByNameStorage(ctx context.Context, name string) (models.UserEntity, error) {
	op := "storage.GetUserByNameStorage"
	log := s.log.With("op", op)

	var user = models.UserEntity{}

	query := `SELECT id, name, email, role_id FROM "users" WHERE name = $1`
	err := pgxscan.Get(ctx, s.Db, &user, query, name)
	if err != nil {
		log.Error(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return user, errorsApp.ErrUserNotFound.Error
		}
		return user, errorsApp.ErrInternalError.Error
	}

	return user, nil
}
