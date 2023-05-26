package users

// UserSignIn is a struct that contains user sign in data.
type RegisterUser struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}
