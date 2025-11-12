package models

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	// ID is stored as a UUID in the database. Use string here so Bun
	// doesn't try to scan it into an integer.
	ID        string    `bun:",pk" json:"id"`
	Name      string    `bun:",notnull" json:"name"`
	Email     string    `bun:",unique,notnull" json:"email"`
	Password  string    `bun:",notnull" json:"-"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
}
