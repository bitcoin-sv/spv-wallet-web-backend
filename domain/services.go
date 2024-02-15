package domain

import (
	db_users "github.com/bitcoin-sv/spv-wallet-web-backend/data/users"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/transactions"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/client"

	"github.com/rs/zerolog"
)

// Services is a struct that contains all services.
type Services struct {
	UsersService        *users.UserService
	TransactionsService *transactions.TransactionService
	WalletClientFactory users.WalletClientFactory
}

// NewServices creates services instance.
func NewServices(usersRepo *db_users.UsersRepository, log *zerolog.Logger) (*Services, error) {
	walletClientFactory := client.NewWalletClientFactory(log)
	adminWalletClient, err := walletClientFactory.CreateAdminClient()
	if err != nil {
		return nil, err
	}

	// Create User services.
	uService := users.NewUserService(usersRepo, adminWalletClient, walletClientFactory, log)

	return &Services{
		UsersService:        uService,
		TransactionsService: transactions.NewTransactionService(adminWalletClient, walletClientFactory, log),
		WalletClientFactory: walletClientFactory,
	}, nil
}
