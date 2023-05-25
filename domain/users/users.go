package users

import (
	"bux-wallet/data/users"
	"time"
)

// User is a struct that contains user data.
type User struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Mnemonic  string    `json:"mnemonic"`
	Xpriv     string    `json:"xpriv"`
	CreatedAt time.Time `json:"created_at"`
}

// toUserDto converts User to UserDto.
func (user *User) toUserDto() *users.UserDto {
	return &users.UserDto{
		Username:  user.Username,
		Password:  user.Password,
		Mnemonic:  user.Mnemonic,
		Xpriv:     user.Xpriv,
		CreatedAt: user.CreatedAt,
	}
}
