package users

import (
	"time"
)

// UserDto is a struct that represent user database record.
type UserDto struct {
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Mnemonic  string    `db:"mnemonic"`
	Xpriv     string    `db:"xpriv"`
	CreatedAt time.Time `db:"created_at"`
}
