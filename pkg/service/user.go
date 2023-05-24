package service

import (
	"bux-wallet/pkg/domain"
	"bux-wallet/pkg/logger"
	"bux-wallet/pkg/repository"

	"github.com/rs/zerolog"
)

type UsersService struct {
	logger zerolog.Logger
	repo   repository.User
}

func NewUsersService(repo repository.User) *UsersService {
	return &UsersService{
		logger: logger.NewServiceLogger(),
		repo:   repo,
	}
}

func (s UsersService) SignUp(input UserSignUpInput) (domain.User, error) {
	return s.repo.Create()
}

func (s UsersService) SignIn(input UserSignInInput) (domain.User, error) {
	return s.repo.GetByCredentials(input.Email, input.Password)
}
