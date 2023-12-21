package buxclient

import (
	"bux-wallet/config"
	"bux-wallet/domain/users"
	"github.com/rs/zerolog"

	"github.com/BuxOrg/go-buxclient"
	"github.com/spf13/viper"
)

type buxclientFactory struct {
	log *zerolog.Logger
}

// NewBuxClientFactory creates instance of Bux Client Factory.
func NewBuxClientFactory(log *zerolog.Logger) users.BuxClientFactory {
	buxClientLogger := log.With().Str("service", "bux-client").Logger()
	return &buxclientFactory{
		log: &buxClientLogger,
	}
}

// CreateAdminBuxClient creates instance of Bux Client with admin keys.
func (bf *buxclientFactory) CreateAdminBuxClient() (users.AdmBuxClient, error) {
	// Get env variables.
	xpriv := viper.GetString(config.EnvBuxAdminXpriv)
	serverUrl, debug, signRequest := getServerData()

	// Init bux client.
	buxClient, err := buxclient.New(
		buxclient.WithXPriv(xpriv),
		buxclient.WithAdminKey(xpriv),
		buxclient.WithHTTP(serverUrl),
		buxclient.WithDebugging(debug),
		buxclient.WithSignRequest(signRequest),
	)

	if err != nil {
		return nil, err
	}

	return &AdminBuxClient{
		client: buxClient,
		log:    bf.log,
	}, nil
}

// CreateWithXpriv creates instance of Bux Client with given xpriv.
func (bf *buxclientFactory) CreateWithXpriv(xpriv string) (users.UserBuxClient, error) {
	// Get env variables.
	serverUrl, debug, signRequest := getServerData()

	// Init bux client with generated xpub.
	buxClient, err := buxclient.New(
		buxclient.WithXPriv(xpriv),
		buxclient.WithHTTP(serverUrl),
		buxclient.WithDebugging(debug),
		buxclient.WithSignRequest(signRequest),
	)

	if err != nil {
		return nil, err
	}

	return &BuxClient{
		client: buxClient,
		log:    bf.log,
	}, nil
}

// CreateWithXpriv creates instance of Bux Client with given xpriv.
func (bf *buxclientFactory) CreateWithAccessKey(accessKey string) (users.UserBuxClient, error) {
	// Get env variables.
	serverUrl, debug, signRequest := getServerData()

	// Init bux client with generated xpub.
	buxClient, err := buxclient.New(
		buxclient.WithAccessKey(accessKey),
		buxclient.WithHTTP(serverUrl),
		buxclient.WithDebugging(debug),
		buxclient.WithSignRequest(signRequest),
	)

	if err != nil {
		return nil, err
	}

	return &BuxClient{
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
