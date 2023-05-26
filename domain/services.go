package domain

import (
	db_users "bux-wallet/data/users"
	"bux-wallet/domain/users"
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
