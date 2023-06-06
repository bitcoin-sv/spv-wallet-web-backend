package buxclient

import (
	"bux-wallet/config"
	"bux-wallet/logging"
	"context"
	"fmt"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/go-buxclient"
	"github.com/libsv/go-bk/bip32"
	"github.com/spf13/viper"
)

// AdminBuxClient is a wrapper for Admin Bux Client.
type AdminBuxClient struct {
	client *buxclient.BuxClient
	log    logging.Logger
}

// RegisterXpub registers xpub in bux.
func (c *AdminBuxClient) RegisterXpub(xpriv *bip32.ExtendedKey) (string, error) {
	// Get xpub from xpriv.
	xpub, err := xpriv.Neuter()

	if err != nil {
		return "", err
	}

	// Register new xpub in BUX.
	err = c.client.NewXpub(
		context.Background(), xpub.String(), &bux.Metadata{},
	)

	if err != nil {
		c.log.Error(err.Error())
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

	// Register new xpub in BUX.
	err := c.client.NewPaymail(context.Background(), xpub, address, alias, alias, &bux.Metadata{})

	if err != nil {
		return "", err
	}
	return address, nil
}
