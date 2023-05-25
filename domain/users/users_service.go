package users

import (
	"bux-wallet/data/users"
	"context"
	"fmt"
)

// UserService represents User service and provide access to repository.
type UserService struct {
	repo UsersRepository
}

// NewUserService creates UserService instance.
func NewUserService(repo *users.UsersRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// InsertUser inserts user to database.
func (s *UserService) InsertUser(user *User) error {
	err := s.repo.InsertUser(context.Background(), user.toUserDto())
	fmt.Println(err)
	return err
}
