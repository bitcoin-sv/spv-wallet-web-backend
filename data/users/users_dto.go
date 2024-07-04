package users

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
)

// UserDto is a struct that represent user database record.
type UserDto struct {
	ID        int       `db:"id"`
	Email     string    `db:"email"`
	Xpriv     string    `db:"xpriv"`
	Paymail   string    `db:"paymail"`
	CreatedAt time.Time `db:"created_at"`
}

// toUser converts UserDto to User.
func (user *UserDto) toUser() *users.User {
	return &users.User{
		ID:        user.ID,
		Email:     user.Email,
		Xpriv:     user.Xpriv,
		Paymail:   user.Paymail,
		CreatedAt: user.CreatedAt,
	}
}
