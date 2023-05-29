package buxclient

import (
	"bux-wallet/config"

	"github.com/BuxOrg/go-buxclient"
	"github.com/spf13/viper"
)

// AdminBuxClient is a wrapper for Admin Bux Client.
type AdminBuxClient struct {
	client *buxclient.BuxClient
}

// BuxClient is a wrapper for Bux Client.
type BuxClient struct {
	client *buxclient.BuxClient
}

// CreateAdminBuxClient creates instance of Bux Client with admin keys.
func CreateAdminBuxClient() (*AdminBuxClient, error) {
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

	return &AdminBuxClient{client: buxClient}, nil
}

// CreateBuxClient creates instance of Bux Client with user xpub.
func CreateBuxClient(xpub string) (*BuxClient, error) {
	// Get env variables.
	serverUrl := viper.GetString(config.EnvBuxServerUrl)
	debug := viper.GetBool(config.EnvBuxWithDebug)
	signRequest := viper.GetBool(config.EnvBuxSignRequest)

	// Init bux client.
	buxClient, err := buxclient.New(
		buxclient.WithXPub(xpub),
		buxclient.WithHTTP(serverUrl),
		buxclient.WithDebugging(debug),
		buxclient.WithSignRequest(signRequest),
	)

	if err != nil {
		return nil, err
	}

	return &BuxClient{buxClient}, nil
}
