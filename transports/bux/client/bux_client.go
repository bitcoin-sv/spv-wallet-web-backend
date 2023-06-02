package buxclient

import (
	"bux-wallet/config"
	"bux-wallet/logging"

	"github.com/BuxOrg/go-buxclient"
	"github.com/spf13/viper"
)

// AdminBuxClient is a wrapper for Admin Bux Client.
type AdminBuxClient struct {
	client *buxclient.BuxClient
	log    logging.Logger
}

// BuxClient is a wrapper for Bux Client.
type BuxClient struct {
	client *buxclient.BuxClient
	log    logging.Logger
}

// CreateAdminBuxClient creates instance of Bux Client with admin keys.
func CreateAdminBuxClient(lf logging.LoggerFactory) (*AdminBuxClient, error) {
	// Get env variables.
	xpriv := viper.GetString(config.EnvBuxAdminXpriv)
	serverUrl := viper.GetString(config.EnvBuxServerUrl)
	debug := viper.GetBool(config.EnvBuxWithDebug)
	signRequest := viper.GetBool(config.EnvBuxSignRequest)

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
		log:    lf.NewLogger("admin-bux-client"),
	}, nil
}

// CreateBuxClientFromRawXpriv creates instance of Bux Client with user raw xpriv.
func CreateBuxClientFromRawXpriv(rawXpriv string) (*BuxClient, error) {
	// Get env variables.
	serverUrl := viper.GetString(config.EnvBuxServerUrl)
	debug := viper.GetBool(config.EnvBuxWithDebug)
	signRequest := viper.GetBool(config.EnvBuxSignRequest)

	// Init bux client with generated xpub.
	buxClient, err := buxclient.New(
		buxclient.WithXPriv(rawXpriv),
		buxclient.WithHTTP(serverUrl),
		buxclient.WithDebugging(debug),
		buxclient.WithSignRequest(signRequest),
	)

	if err != nil {
		return nil, err
	}

	lf := logging.DefaultLoggerFactory()

	return &BuxClient{
		client: buxClient,
		log:    lf.NewLogger("bux-client"),
	}, nil
}
