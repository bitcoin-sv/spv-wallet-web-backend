package domain

import (
	db_users "bux-wallet/data/users"
	"bux-wallet/domain/transactions"
	"bux-wallet/domain/users"
	"bux-wallet/logging"
	buxclient "bux-wallet/transports/bux/client"
)

// Services is a struct that contains all services.
type Services struct {
	UsersService        *users.UserService
	TransactionsService *transactions.TransactionService
	BuxClientFactory    users.BuxClientFactory
}

// NewServices creates services instance.
func NewServices(usersRepo *db_users.UsersRepository, lf logging.LoggerFactory) (*Services, error) {
	bf := buxclient.NewBuxClientFactory(lf)
	adminBuxClient, err := bf.CreateAdminBuxClient()
	if err != nil {
		return nil, err
	}

	// Create User services.
	uService := users.NewUserService(usersRepo, adminBuxClient, bf, lf)

	return &Services{
		UsersService:        uService,
		TransactionsService: transactions.NewTransactionService(adminBuxClient, bf, lf),
		BuxClientFactory:    bf,
	}, nil
}
