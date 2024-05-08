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
	xpriv := viper.GetString(config.EnvAdminXpriv)
	serverUrl, signRequest := getServerData()

	adminWalletClient, err := walletclient.New(
		walletclient.WithXPriv(xpriv),
		walletclient.WithAdminKey(xpriv),
		walletclient.WithHTTP(serverUrl),
		walletclient.WithSignRequest(signRequest),
	)

	if err != nil {
		return nil, err
	}

	return &AdminWalletClient{
		client: adminWalletClient,
		log:    bf.log,
	}, nil
}

// CreateWithXpriv returns UserWalletClient as spv-wallet-go-client instance with given xpriv.
func (bf *walletClientFactory) CreateWithXpriv(xpriv string) (users.UserWalletClient, error) {
	serverUrl, signRequest := getServerData()

	userWalletClient, err := walletclient.New(
		walletclient.WithXPriv(xpriv),
		walletclient.WithHTTP(serverUrl),
		walletclient.WithSignRequest(signRequest),
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		client: userWalletClient,
		log:    bf.log,
	}, nil
}

// CreateWithXpriv returns UserWalletClient as spv-wallet-go-client instance with given access key.
func (bf *walletClientFactory) CreateWithAccessKey(accessKey string) (users.UserWalletClient, error) {
	// Get env variables.
	serverUrl, signRequest := getServerData()

	userWalletClient, err := walletclient.New(
		walletclient.WithAccessKey(accessKey),
		walletclient.WithHTTP(serverUrl),
		walletclient.WithSignRequest(signRequest),
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		client: userWalletClient,
		log:    bf.log,
	}, nil
}

func getServerData() (serverUrl string, signRequest bool) {
	// Get env variables.
	serverUrl = viper.GetString(config.EnvServerUrl)
	signRequest = viper.GetBool(config.EnvSignRequest)

	return serverUrl, signRequest
}
