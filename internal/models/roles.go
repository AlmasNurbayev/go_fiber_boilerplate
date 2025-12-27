package models

import "time"

type RoleEntity struct {
	Id           int64     `db:"id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	Changed_date time.Time `json:"changed_date" db:"changed_date"`
	Create_date  time.Time `json:"create_date" db:"create_date"`
}
