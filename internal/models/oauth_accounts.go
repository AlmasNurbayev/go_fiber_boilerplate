package models

import (
	"time"
)

type OauthAccountEntity struct {
	Id               int64     `db:"id"`
	User_id          int64     `db:"user_id"`
	Provider         string    `db:"provider"`
	Provider_user_id string    `db:"provider_user_id"`
	Changed_date     time.Time `db:"changed_date"`
	Create_date      time.Time `db:"create_date"`
}
