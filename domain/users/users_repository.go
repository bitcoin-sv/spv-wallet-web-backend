package users

import (
	"bux-wallet/data/users"
	"context"
)

// UsersRepository is an interface which defines methods for UsersRepository.
type UsersRepository interface {
	InsertUser(ctx context.Context, user *users.UserDto) error
	GetUserByEmail(ctx context.Context, email string) (*users.UserDto, error)
}
