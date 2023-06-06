package users

import (
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

// AuthenticatedUser is a struct that contains authenticated user data.
type AuthenticatedUser struct {
	User      *User
	AccessKey AccessKey
}

// AccessKey is a struct that contains access key data.
type AccessKey struct {
	Id  string `json:"id"`
	Key string `json:"key"`
}
