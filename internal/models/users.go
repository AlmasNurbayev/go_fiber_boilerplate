package models

import (
	"time"

	"github.com/guregu/null/v6"
)

type UserEntity struct {
	Id            int64       `db:"id"`
	Phone_number  null.String `db:"phone_number"`
	Email         null.String `db:"email"`
	Name          string      `db:"name"`
	Password_hash string      `db:"password_hash"`
	Role_id       int64       `db:"role_id"`
	Changed_date  time.Time   `db:"changed_date"`
	Create_date   time.Time   `db:"create_date"`
}
