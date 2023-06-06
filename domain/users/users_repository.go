package users

import (
	"context"
)

// UsersRepository is an interface which defines methods for UsersRepository.
type UsersRepository interface {
	InsertUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserById(ctx context.Context, id int) (*User, error)
}
