package users

import (
	"time"
)

// UserDto is a struct that represent user database record.
type UserDto struct {
	Email     string    `db:"email"`
	Xpriv     string    `db:"xpriv"`
	Paymail   string    `db:"paymail"`
	CreatedAt time.Time `db:"created_at"`
}
