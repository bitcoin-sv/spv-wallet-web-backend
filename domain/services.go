package domain

import (
	db_users "bux-wallet/data/users"
	"bux-wallet/domain/users"
	"bux-wallet/logging"
	buxclient "bux-wallet/transports/bux/client"
)

// Services is a struct that contains all services.
type Services struct {
	UsersService     *users.UserService
	BuxClientFactory users.BuxClientFactory
}

// NewServices creates services instance.
func NewServices(usersRepo *db_users.UsersRepository, lf logging.LoggerFactory) (*Services, error) {
	bf := buxclient.NewBuxClientFactory(lf)

	// Create User services.
	uService, err := users.NewUserService(usersRepo, bf, lf)
	if err != nil {
		return nil, err
	}

	return &Services{
		UsersService:     uService,
		BuxClientFactory: bf,
	}, nil
}
