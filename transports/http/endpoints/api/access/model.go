package access

// SignInUser is a struct that contains user sign in data.
type SignInUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
