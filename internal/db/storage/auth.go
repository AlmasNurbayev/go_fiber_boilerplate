package storage

import (
	"context"
	"errors"
	"log/slog"
	"strconv"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) NewUser(ctx context.Context, user models.UserEntity) (models.UserEntity, *errorsApp.DbError) {
	op := "storage.NewUser"
	log := s.log.With("op", op)

	// искусственное замедление запроса
	// var temp string
	// err2 := pgxscan.Get(ctx, s.db, &temp, "SELECT pg_sleep(18)")
	// if err2 != nil {
	// 	s.log.Error("canceled query DB", "error", err2)
	// 	return user, err2
	// }

	query := `INSERT INTO "users" (name, phone_number, email, password_hash, role_id) VALUES ($1, $2, $3, $4, $5) RETURNING *`

	rows, err := s.Db.Query(ctx, query, user.Name, user.Phone_number, user.Email, user.Password_hash, user.Role_id)
	if err != nil {
		log.Error(err.Error())
		return user, mapPgError(err)
	}

	savedUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.UserEntity])
	if err != nil {
		log.Error(err.Error())
		return user, mapPgError(err)
	}

	return savedUser, nil
}

func (s *Storage) GetRoleById(ctx context.Context, id int64) (models.RoleEntity, *errorsApp.DbError) {
	op := "storage.GetRoleById"
	log := s.log.With("op", op)

	query := `SELECT * FROM "roles" WHERE id = $1`

	rows, err := s.Db.Query(ctx, query, id)
	if err != nil {
		log.Error(err.Error())
		return models.RoleEntity{}, mapPgError(err)
	}

	role, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.RoleEntity])
	if err != nil {
		log.Error(err.Error())
		return models.RoleEntity{}, mapPgError(err)
	}

	return role, nil
}

func (s *Storage) GetUserById(ctx context.Context, id int64) (models.UserEntity, *errorsApp.DbError) {
	op := "storage.GetUserByIdStorage"
	log := s.log.With("op", op)
	var user = []models.UserEntity{}

	query := `SELECT * FROM "users" WHERE id = $1`

	err := pgxscan.Get(ctx, s.Db, &user, query, id)
	if err != nil {
		log.Error(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return user[0], &errorsApp.DbError{
				Field:   "id",
				Data:    id,
				Message: "user not found",
				Error:   errors.New("user with id " + strconv.FormatInt(id, 10) + " not found"),
			}
		}
		return user[0], mapPgError(err)
	}
	return user[0], nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (models.UserEntity, *errorsApp.DbError) {
	op := "storage.GetUserByEmail"
	log := s.log.With("op", op)
	var users = models.UserEntity{}

	log.Debug("looking for user by email", slog.String("email", email))

	query := `SELECT * FROM "users" WHERE email = $1`

	err := pgxscan.Get(ctx, s.Db, &users, query, email)
	if err != nil {
		log.Error(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return users, &errorsApp.DbError{
				Type:    "not_found",
				Field:   "email",
				Data:    email,
				Message: "user not found",
				Error:   errors.New("user with email " + email + " not found"),
			}
		}
		return users, mapPgError(err)
	}
	return users, nil
}

func (s *Storage) GetUserByPhoneNumber(ctx context.Context, phone_number string) (models.UserEntity, *errorsApp.DbError) {
	op := "storage.GetUserByPhoneNumber"
	log := s.log.With("op", op)
	var user = models.UserEntity{}

	query := `SELECT * FROM "users" WHERE phone_number = $1`

	err := pgxscan.Get(ctx, s.Db, &user, query, phone_number)
	if err != nil {
		log.Error(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return user, &errorsApp.DbError{
				Field:   "phone_number",
				Data:    phone_number,
				Message: "user not found",
				Error:   errors.New("user with phone_number " + phone_number + " not found"),
			}
		}
		return user, mapPgError(err)
	}
	return user, nil
}
