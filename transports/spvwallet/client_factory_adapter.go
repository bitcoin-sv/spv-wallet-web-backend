package spvwallet

import (
	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type walletClientFactory struct {
	log *zerolog.Logger
}

// NewWalletClientFactory implements the ClientFactory.
func NewWalletClientFactory(log *zerolog.Logger) users.WalletClientFactory {
	logger := log.With().Str("service", "spv-wallet-client").Logger()
	return &walletClientFactory{
		log: &logger,
	}
}

// CreateAdminClient returns AdminWalletClient as spv-wallet-go-client instance with admin key.
func (bf *walletClientFactory) CreateAdminClient() (users.AdminWalletClient, error) {
	return NewAdminClientAdapter(bf.log)
}

// CreateWithXpriv returns UserWalletClient as spv-wallet-go-client instance with given xpriv.
func (bf *walletClientFactory) CreateWithXpriv(xpriv string) (users.UserWalletClient, error) {
	return NewUserClientAdapterWithXPriv(bf.log, xpriv)
}

// CreateWithAccessKey returns UserWalletClient as spv-wallet-go-client instance with given access key.
func (bf *walletClientFactory) CreateWithAccessKey(accessKey string) (users.UserWalletClient, error) {
	return NewUserClientAdapterWithAccessKey(bf.log, accessKey)
}

func getServerData() string {
	// Get env variables.
	return viper.GetString(config.EnvServerURL)
}
