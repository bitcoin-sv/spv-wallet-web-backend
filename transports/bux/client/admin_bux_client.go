package buxclient

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"

	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/go-buxclient"
	"github.com/libsv/go-bk/bip32"
	"github.com/spf13/viper"

	"bux-wallet/config"
)

// AdminBuxClient is a wrapper for Admin Bux Client.
type AdminBuxClient struct {
	client *buxclient.BuxClient
	log    *zerolog.Logger
}

// RegisterXpub registers xpub in bux.
func (c *AdminBuxClient) RegisterXpub(xpriv *bip32.ExtendedKey) (string, error) {
	// Get xpub from xpriv.
	xpub, err := xpriv.Neuter()

	if err != nil {
		c.log.Error().Msgf("Error while creating new xPub: %v", err.Error())
		return "", err
	}

	// Register new xpub in BUX.
	err = c.client.NewXpub(
		context.Background(), xpub.String(), &buxmodels.Metadata{},
	)

	if err != nil {
		c.log.Error().
			Str("xpub", xpub.String()).
			Msgf("Error while registering new xPub: %v", err.Error())
		return "", err
	}

	return xpub.String(), nil
}

// RegisterPaymail registers new paymail in bux.
func (c *AdminBuxClient) RegisterPaymail(alias, xpub string) (string, error) {
	// Get paymail domain from env.
	domain := viper.GetString(config.EnvBuxPaymailDomain)

	// Create paymail address.
	address := fmt.Sprintf("%s@%s", alias, domain)

	// Get avatar url from env.
	avatar := viper.GetString(config.EnvBuxPaymailAvatar)

	// Register new xpub in BUX.
	err := c.client.NewPaymail(context.Background(), xpub, address, avatar, alias, &buxmodels.Metadata{})

	if err != nil {
		c.log.Error().Msgf("Error while registering new paymail: %v", err.Error())
		return "", err
	}
	return address, nil
}
