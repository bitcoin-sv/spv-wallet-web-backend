package client

import (
	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"

	"github.com/rs/zerolog"

	walletclient "github.com/BuxOrg/go-buxclient"
	"github.com/spf13/viper"
)

type clientFactory struct {
	log *zerolog.Logger
}

// NewClientFactory implements the ClientFactory.
func NewClientFactory(log *zerolog.Logger) users.ClientFactory {
	logger := log.With().Str("service", "spv-wallet-client").Logger()
	return &clientFactory{
		log: &logger,
	}
}

// CreateAdminClient returns AdminClient as spv-wallet-go-client instance with admin key.
func (bf *clientFactory) CreateAdminClient() (users.AdminClient, error) {
	xpriv := viper.GetString(config.EnvAdminXpriv)
	serverUrl, debug, signRequest := getServerData()

	adminClient, err := walletclient.New(
		walletclient.WithXPriv(xpriv),
		walletclient.WithAdminKey(xpriv),
		walletclient.WithHTTP(serverUrl),
		walletclient.WithDebugging(debug),
		walletclient.WithSignRequest(signRequest),
	)

	if err != nil {
		return nil, err
	}

	return &AdminClient{
		client: adminClient,
		log:    bf.log,
	}, nil
}

// CreateWithXpriv returns UserClient as spv-wallet-go-client instance with given xpriv.
func (bf *clientFactory) CreateWithXpriv(xpriv string) (users.UserClient, error) {
	serverUrl, debug, signRequest := getServerData()

	userClient, err := walletclient.New(
		walletclient.WithXPriv(xpriv),
		walletclient.WithHTTP(serverUrl),
		walletclient.WithDebugging(debug),
		walletclient.WithSignRequest(signRequest),
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		client: userClient,
		log:    bf.log,
	}, nil
}

// CreateWithXpriv returns UserClient as spv-wallet-go-client instance with given access key.
func (bf *clientFactory) CreateWithAccessKey(accessKey string) (users.UserClient, error) {
	// Get env variables.
	serverUrl, debug, signRequest := getServerData()

	userclient, err := walletclient.New(
		walletclient.WithAccessKey(accessKey),
		walletclient.WithHTTP(serverUrl),
		walletclient.WithDebugging(debug),
		walletclient.WithSignRequest(signRequest),
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		client: userclient,
		log:    bf.log,
	}, nil
}

func getServerData() (serverUrl string, debug, signRequest bool) {
	// Get env variables.
	serverUrl = viper.GetString(config.EnvServerUrl)
	debug = viper.GetBool(config.EnvWithDebug)
	signRequest = viper.GetBool(config.EnvSignRequest)

	return serverUrl, debug, signRequest
}
