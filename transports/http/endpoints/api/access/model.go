package access

import "web-backend/domain/users"

// SignInUser is a struct that contains user sign in data.
type SignInUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignInResponse is a struct that represents struct sended after user sign in.
type SignInResponse struct {
	Paymail string        `json:"paymail"`
	Balance users.Balance `json:"balance"`
}
