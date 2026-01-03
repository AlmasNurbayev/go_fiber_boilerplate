package storage

import (
	"context"
	"errors"
	"strconv"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib/errorsApp"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) NewOauthAccount(ctx context.Context, user models.OauthAccountEntity) (models.OauthAccountEntity, *errorsApp.DbError) {
	op := "storage.NewOauthAccount"
	log := s.log.With("op", op)

	query := `INSERT INTO "oauth_accounts" (user_id, provider, provider_user_id) VALUES ($1, $2, $3 RETURNING *`

	rows, err := s.Db.Query(ctx, query, user.User_id, user.Provider, user.Provider_user_id)
	if err != nil {
		log.Error(err.Error())
		return user, mapPgError(err)
	}

	savedOauth, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.OauthAccountEntity])
	if err != nil {
		log.Error(err.Error())
		return user, mapPgError(err)
	}

	return savedOauth, nil
}

func (s *Storage) GetOauthAccountById(ctx context.Context, id int64) (models.OauthAccountEntity, *errorsApp.DbError) {
	op := "storage.GetOauthAccountById"
	log := s.log.With("op", op)

	query := `SELECT * FROM "oauth_accounts" WHERE id = $1`
	account := models.OauthAccountEntity{}

	err := pgxscan.Get(ctx, s.Db, &account, query, id)
	if err != nil {
		log.Error(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return account, &errorsApp.DbError{
				Type:    "not_found",
				Field:   "id",
				Data:    id,
				Message: "oauth account not found",
				Error:   errors.New("oauth account with id " + strconv.FormatInt(id, 10) + " not found"),
			}
		}
		return account, mapPgError(err)
	}
	return account, nil
}

func (s *Storage) GetOauthAccountByUserId(ctx context.Context, id int64) (models.OauthAccountEntity, *errorsApp.DbError) {
	op := "storage.GetOauthAccountByUserId"
	log := s.log.With("op", op)

	query := `SELECT * FROM "oauth_accounts" WHERE user_id = $1`
	account := models.OauthAccountEntity{}

	err := pgxscan.Get(ctx, s.Db, &account, query, id)
	if err != nil {
		log.Error(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return account, &errorsApp.DbError{
				Type:    "not_found",
				Field:   "user_id",
				Data:    id,
				Message: "oauth account not found",
				Error:   errors.New("oauth account with user_id " + strconv.FormatInt(id, 10) + " not found"),
			}
		}
		return account, mapPgError(err)
	}
	return account, nil
}
