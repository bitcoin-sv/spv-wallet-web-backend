package service

import (
	"bux-wallet/pkg/domain"
	"bux-wallet/pkg/repository"
)

type UserSignUpInput struct {
	Name     string
	Email    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type Users interface {
	SignUp(input UserSignUpInput) (domain.User, error)
	SignIn(input UserSignInInput) (domain.User, error)
}

type Services struct {
	Users Users
}

func NewServices(r *repository.Repositories) *Services {
	return &Services{
		Users: NewUsersService(r.User),
	}

}
