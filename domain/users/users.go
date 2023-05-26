package users

import (
	"bux-wallet/data/users"
	"time"
)

// User is a struct that contains user data.
type User struct {
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Mnemonic  string    `json:"mnemonic"`
	Xpriv     string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// toUserDto converts User to UserDto.
func (user *User) toUserDto() *users.UserDto {
	return &users.UserDto{
		Email:     user.Email,
		Password:  user.Password,
		Mnemonic:  user.Mnemonic,
		Xpriv:     user.Xpriv,
		CreatedAt: user.CreatedAt,
	}
}
