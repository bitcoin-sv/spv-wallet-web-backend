package spvwallet

import (
	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet-web-backend/spverrors"
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
	serverURL := getServerData()

	adminWalletClient, err := walletclient.NewWithAdminKey(serverURL, adminKey)
	if err != nil {
		return nil, spverrors.ErrCreateClientAdminKey.Wrap(err)
	}

	return &AdminWalletClient{
		client: adminWalletClient,
		log:    bf.log,
	}, nil
}

// CreateWithXpriv returns UserWalletClient as spv-wallet-go-client instance with given xpriv.
func (bf *walletClientFactory) CreateWithXpriv(xpriv string) (users.UserWalletClient, error) {
	serverURL := getServerData()

	userWalletClient, err := walletclient.NewWithXPriv(serverURL, xpriv)
	if err != nil {
		return nil, spverrors.ErrInvalidCredentials.Wrap(err)
	}

	return &Client{
		client: userWalletClient,
		log:    bf.log,
	}, nil
}

// CreateWithAccessKey returns UserWalletClient as spv-wallet-go-client instance with given access key.
func (bf *walletClientFactory) CreateWithAccessKey(accessKey string) (users.UserWalletClient, error) {
	serverURL := getServerData()

	userWalletClient, err := walletclient.NewWithAccessKey(serverURL, accessKey)
	if err != nil {
		return nil, spverrors.ErrInvalidCredentials.Wrap(err)
	}

	return &Client{
		client: userWalletClient,
		log:    bf.log,
	}, nil
}

func getServerData() string {
	// Get env variables.
	return viper.GetString(config.EnvServerURL)
}
