package domain

import (
	"bux-wallet/domain/users"
	db_users "bux-wallet/data/users"
)

// Services is a struct that contains all services.
type Services struct {
	UsersService *users.UserService
}

// NewServices creates services instance.
func NewServices(usersRepo *db_users.UsersRepository) *Services {
	return &Services{
		UsersService: users.NewUserService(usersRepo),
	}
}