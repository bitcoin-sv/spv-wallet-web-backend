package users

// RegisterUser is a struct that contains user register data.
type RegisterUser struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

// RegisterResponse represents response that is sent after user creation.
type RegisterResponse struct {
	Mnemonic string `json:"mnemonic"`
	Paymail  string `json:"paymail"`
}

type UserResponse struct {
	UserId  int    `json:"userId"`
	Token   string `json:"token"`
	Paymail string `json:"paymail"`
}
