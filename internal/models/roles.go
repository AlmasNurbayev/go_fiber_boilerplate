package models

import "time"

type RoleEntity struct {
	Id           int64     `db:"id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	Changed_date time.Time `db:"changed_date"`
	Create_date  time.Time `db:"create_date"`
}
