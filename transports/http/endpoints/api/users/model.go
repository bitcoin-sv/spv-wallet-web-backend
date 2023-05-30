package users

// RegisterUser is a struct that contains user register data.
type RegisterUser struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

// RegisterReposne is a struct that contains user mnemonic.
type RegisterReposne struct {
	Mnemonic string `json:"mnemonic"`
	Paymail  string `json:"paymail"`
}
