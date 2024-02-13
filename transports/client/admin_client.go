package client

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	models "github.com/BuxOrg/bux-models"
	walletclient "github.com/BuxOrg/go-buxclient"
	"github.com/libsv/go-bk/bip32"
	"github.com/spf13/viper"

	"bux-wallet/config"
)

// AdminClient is a wrapper for Admin SPV Wallet Client.
type AdminClient struct {
	client *walletclient.BuxClient
	log    *zerolog.Logger
}

// RegisterXpub registers xpub in SPV Wallet.
func (c *AdminClient) RegisterXpub(xpriv *bip32.ExtendedKey) (string, error) {
	// Get xpub from xpriv.
	xpub, err := xpriv.Neuter()

	if err != nil {
		c.log.Error().Msgf("Error while creating new xPub: %v", err.Error())
		return "", err
	}

	// Register new xpub in SPV Wallet.
	err = c.client.NewXpub(
		context.Background(), xpub.String(), &models.Metadata{},
	)

	if err != nil {
		c.log.Error().
			Str("xpub", xpub.String()).
			Msgf("Error while registering new xPub: %v", err.Error())
		return "", err
	}

	return xpub.String(), nil
}

// RegisterPaymail registers new paymail in SPV Wallet.
func (c *AdminClient) RegisterPaymail(alias, xpub string) (string, error) {
	// Get paymail domain from env.
	domain := viper.GetString(config.EnvBuxPaymailDomain)

	// Create paymail address.
	address := fmt.Sprintf("%s@%s", alias, domain)

	// Get avatar url from env.
	avatar := viper.GetString(config.EnvBuxPaymailAvatar)

	// Register new xpub in SPV Wallet.
	err := c.client.NewPaymail(context.Background(), xpub, address, avatar, alias, &models.Metadata{})

	if err != nil {
		c.log.Error().Msgf("Error while registering new paymail: %v", err.Error())
		return "", err
	}
	return address, nil
}
