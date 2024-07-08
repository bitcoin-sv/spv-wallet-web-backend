package users

import (
	"time"
)

// User is a struct that contains user data.
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Xpriv     string    `json:"-"` // xPriv encrypted with user password
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
	Balance   Balance
	Xpriv     string `json:"-"` // xPriv should not be exposed to the client
}

// AccessKey is a struct that contains access key data.
type AccessKey struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

// ExchangeRate is a struct that contains exchange rate data.
type ExchangeRate struct {
	Rate float64
}

// Balance is a struct that contains user balance data.
type Balance struct {
	Usd      float64 `json:"usd"`
	Bsv      float64 `json:"bsv"`
	Satoshis uint64  `json:"satoshis"`
}
