package models

import (
	"time"

	"github.com/guregu/null/v6"
)

type UserEntity struct {
	Id            int64       `json:"id" db:"id"`
	Phone_number  string      `json:"phone_number" db:"phone_number"`
	Email         string      `json:"email" db:"email"`
	Name          null.String `json:"name" db:"name"`
	Password_hash null.String `json:"password_hash" db:"password_hash"`
	Role_id       null.Int64  `json:"role" db:"role_id"`
	Changed_date  time.Time   `json:"changed_date" db:"changed_date"`
	Create_date   time.Time   `json:"create_date" db:"create_date"`
}
