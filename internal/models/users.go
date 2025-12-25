package models

import (
	"time"

	"github.com/guregu/null/v6"
)

type UserEntity struct {
	Id            int64       `json:"id" db:"id"`
	Phone_number  null.String `json:"phone_number" db:"phone_number"`
	Email         null.String `json:"email" db:"email"`
	Name          string      `json:"name" db:"name"`
	Password_hash string      `json:"password_hash" db:"password_hash"`
	Role_id       int64       `json:"role" db:"role_id"`
	Changed_date  time.Time   `json:"changed_date" db:"changed_date"`
	Create_date   time.Time   `json:"create_date" db:"create_date"`
}
