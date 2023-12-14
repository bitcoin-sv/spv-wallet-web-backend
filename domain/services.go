package domain

import (
	db_users "bux-wallet/data/users"
	"bux-wallet/domain/transactions"
	"bux-wallet/domain/users"
	buxclient "bux-wallet/transports/bux/client"
	"github.com/rs/zerolog"
)

// Services is a struct that contains all services.
type Services struct {
	UsersService        *users.UserService
	TransactionsService *transactions.TransactionService
	BuxClientFactory    users.BuxClientFactory
}

// NewServices creates services instance.
func NewServices(usersRepo *db_users.UsersRepository, log *zerolog.Logger) (*Services, error) {
	bf := buxclient.NewBuxClientFactory(log)
	adminBuxClient, err := bf.CreateAdminBuxClient()
	if err != nil {
		return nil, err
	}

	// Create User services.
	uService := users.NewUserService(usersRepo, adminBuxClient, bf, log)

	return &Services{
		UsersService:        uService,
		TransactionsService: transactions.NewTransactionService(adminBuxClient, bf, log),
		BuxClientFactory:    bf,
	}, nil
}
