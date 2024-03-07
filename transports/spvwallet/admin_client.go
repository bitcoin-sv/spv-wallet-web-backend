package spvwallet

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/libsv/go-bk/bip32"
	"github.com/spf13/viper"

	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
)

// AdminWalletClient is a wrapper for Admin SPV Wallet Client.
type AdminWalletClient struct {
	client *walletclient.WalletClient
	log    *zerolog.Logger
}

// RegisterXpub registers xpub in SPV Wallet.
func (c *AdminWalletClient) RegisterXpub(xpriv *bip32.ExtendedKey) (string, error) {
	// Get xpub from xpriv.
	xpub, err := xpriv.Neuter()

	if err != nil {
		c.log.Error().Msgf("Error while creating new xPub: %v", err.Error())
		return "", err
	}

	// Register new xpub in SPV Wallet.
	err = c.client.AdminNewXpub(
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
func (c *AdminWalletClient) RegisterPaymail(alias, xpub string) (string, error) {
	// Get paymail domain from env.
	domain := viper.GetString(config.EnvPaymailDomain)

	// Create paymail address.
	address := fmt.Sprintf("%s@%s", alias, domain)

	// Get avatar url from env.
	avatar := viper.GetString(config.EnvPaymailAvatar)

	_, err := c.client.AdminCreatePaymail(context.Background(), xpub, address, alias, avatar)

	if err != nil {
		c.log.Error().Msgf("Error while registering new paymail: %v", err.Error())
		return "", err
	}
	return address, nil
}

func (c *AdminWalletClient) GetSharedConfig() (*models.SharedConfig, error) {
	sharedConfig, err := c.client.AdminGetSharedConfig(context.Background())
	if err != nil {
		c.log.Error().Msgf("Error while getting shared config: %v", err.Error())
		return nil, err
	}
	return sharedConfig, nil
}
