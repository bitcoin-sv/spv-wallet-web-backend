package users

import (
	"context"
)

// Repository is an interface which defines methods for Repository.
type Repository interface {
	InsertUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
}
