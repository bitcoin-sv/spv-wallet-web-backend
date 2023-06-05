package users

import (
	"bux-wallet/domain/users"
	"time"
)

// UserDto is a struct that represent user database record.
type UserDto struct {
	Id        int       `db:"id"`
	Email     string    `db:"email"`
	Xpriv     string    `db:"xpriv"`
	Paymail   string    `db:"paymail"`
	CreatedAt time.Time `db:"created_at"`
}

// toUserDto converts User to UserDto.
func toUserDto(user *users.User) *UserDto {
	return &UserDto{
		Id:        user.Id,
		Email:     user.Email,
		Xpriv:     user.Xpriv,
		Paymail:   user.Paymail,
		CreatedAt: user.CreatedAt,
	}
}

// toUser converts UserDto to User.
func (user *UserDto) toUser() *users.User {
	return &users.User{
		Id:        user.Id,
		Email:     user.Email,
		Xpriv:     user.Xpriv,
		Paymail:   user.Paymail,
		CreatedAt: user.CreatedAt,
	}
}
