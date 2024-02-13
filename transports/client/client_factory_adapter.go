package client

import (
	"web-backend/config"
	"web-backend/domain/users"

	"github.com/rs/zerolog"

	walletclient "github.com/BuxOrg/go-buxclient"
	"github.com/spf13/viper"
)

type clientFactory struct {
	log *zerolog.Logger
}

// NewClientFactory creates instance of Bux Client Factory.
func NewClientFactory(log *zerolog.Logger) users.BuxClientFactory {
	buxClientLogger := log.With().Str("service", "spv-wallet-client").Logger()
	return &clientFactory{
		log: &buxClientLogger,
	}
}

// CreateAdminClient creates instance of Bux Client with admin keys.
func (bf *clientFactory) CreateAdminClient() (users.AdmBuxClient, error) {
	// Get env variables.
	xpriv := viper.GetString(config.EnvBuxAdminXpriv)
	serverUrl, debug, signRequest := getServerData()

	// Init bux client.
	buxClient, err := walletclient.New(
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
		client: buxClient,
		log:    bf.log,
	}, nil
}

// CreateWithXpriv creates instance of Bux Client with given xpriv.
func (bf *clientFactory) CreateWithXpriv(xpriv string) (users.UserBuxClient, error) {
	// Get env variables.
	serverUrl, debug, signRequest := getServerData()

	// Init bux client with generated xpub.
	buxClient, err := walletclient.New(
		walletclient.WithXPriv(xpriv),
		walletclient.WithHTTP(serverUrl),
		walletclient.WithDebugging(debug),
		walletclient.WithSignRequest(signRequest),
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		client: buxClient,
		log:    bf.log,
	}, nil
}

// CreateWithXpriv creates instance of Bux Client with given xpriv.
func (bf *clientFactory) CreateWithAccessKey(accessKey string) (users.UserBuxClient, error) {
	// Get env variables.
	serverUrl, debug, signRequest := getServerData()

	// Init bux client with generated xpub.
	buxClient, err := walletclient.New(
		walletclient.WithAccessKey(accessKey),
		walletclient.WithHTTP(serverUrl),
		walletclient.WithDebugging(debug),
		walletclient.WithSignRequest(signRequest),
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		client: buxClient,
		log:    bf.log,
	}, nil
}

func getServerData() (serverUrl string, debug, signRequest bool) {
	// Get env variables.
	serverUrl = viper.GetString(config.EnvBuxServerUrl)
	debug = viper.GetBool(config.EnvBuxWithDebug)
	signRequest = viper.GetBool(config.EnvBuxSignRequest)

	return serverUrl, debug, signRequest
}
