package users

import (
	"bux-wallet/data/users"
	"time"
)

// User is a struct that contains user data.
type User struct {
	Id        int       `json:"id"`
	Email     string    `json:"email"`
	Xpriv     string    `json:"-"`
	Paymail   string    `json:"paymail"`
	CreatedAt time.Time `json:"created_at"`
}

// CreatedUser is a struct that contains new user information used to create http response.
type CreatedUser struct {
	User     *User
	Mnemonic string
}

// SignInUser is a struct that contains user sign in data.
type SignInUser struct {
	User        *User
	AccessKeyId string
}

// toUserDto converts User to UserDto.
func (user *User) toUserDto() *users.UserDto {
	return &users.UserDto{
		Id:        user.Id,
		Email:     user.Email,
		Xpriv:     user.Xpriv,
		Paymail:   user.Paymail,
		CreatedAt: user.CreatedAt,
	}
}

// toUser converts UserDto to User.
func toUser(user *users.UserDto) *User {
	return &User{
		Id:        user.Id,
		Email:     user.Email,
		Xpriv:     user.Xpriv,
		Paymail:   user.Paymail,
		CreatedAt: user.CreatedAt,
	}
}
