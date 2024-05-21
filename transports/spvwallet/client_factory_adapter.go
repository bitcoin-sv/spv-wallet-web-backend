package spvwallet

import (
	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
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
	adminKey := viper.GetString(config.EnvAdminXpriv)
	serverUrl := getServerData()

	adminWalletClient := walletclient.NewWithAdminKey(serverUrl, adminKey)

	return &AdminWalletClient{
		client: adminWalletClient,
		log:    bf.log,
	}, nil
}

// CreateWithXpriv returns UserWalletClient as spv-wallet-go-client instance with given xpriv.
func (bf *walletClientFactory) CreateWithXpriv(xpriv string) (users.UserWalletClient, error) {
	serverUrl := getServerData()

	userWalletClient := walletclient.NewWithXPriv(serverUrl, xpriv)

	return &Client{
		client: userWalletClient,
		log:    bf.log,
	}, nil
}

// CreateWithXpriv returns UserWalletClient as spv-wallet-go-client instance with given access key.
func (bf *walletClientFactory) CreateWithAccessKey(accessKey string) (users.UserWalletClient, error) {
	// Get env variables.
	serverUrl := getServerData()

	userWalletClient := walletclient.NewWithAccessKey(serverUrl, accessKey)

	return &Client{
		client: userWalletClient,
		log:    bf.log,
	}, nil
}

func getServerData() string {
	// Get env variables.
	return viper.GetString(config.EnvServerUrl)
}
