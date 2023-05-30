package users

// RegisterUser is a struct that contains user register data.
type RegisterUser struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

// RegisterResposne is a struct that contains user mnemonic.
type RegisterResposne struct {
	Mnemonic string `json:"mnemonic"`
	Paymail  string `json:"paymail"`
}
