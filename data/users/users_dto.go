package users

import (
	"time"
)

// UserDto is a struct that represent user database record.
type UserDto struct {
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Mnemonic  string    `db:"mnemonic"`
	Xpriv     string    `db:"xpriv"`
	CreatedAt time.Time `db:"created_at"`
}
