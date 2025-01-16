package spvwallet

import (
	"context"
	"fmt"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	walletclientCfg "github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/libsv/go-bk/bip32"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type adminClientAdapter struct {
	log *zerolog.Logger
	api *walletclient.AdminAPI
}

func (a *adminClientAdapter) RegisterXpub(xpriv *bip32.ExtendedKey) (string, error) {
	// Get xpub from xpriv.
	xpub, err := xpriv.Neuter()
	if err != nil {
		a.log.Error().Msgf("Error while returning a new extended public key from the xPriv: %v", err.Error())
		return "", errors.Wrap(err, "error while returning a new extended public key from the xPriv")
	}

	_, err = a.api.CreateXPub(context.Background(), &commands.CreateUserXpub{XPub: xpub.String()})
	if err != nil {
		a.log.Error().Str("xpub", xpub.String()).Msgf("Error while creating new xPub: %v", err.Error())
		return "", errors.Wrap(err, "error while creating new xPub")
	}

	return xpub.String(), nil
}

func (a *adminClientAdapter) RegisterPaymail(alias, xpub string) (string, error) {
	// Get paymail domain from env.
	domain := viper.GetString(config.EnvPaymailDomain)

	// Create paymail address.
	address := fmt.Sprintf("%s@%s", alias, domain)

	// Get avatar url from env.
	avatar := viper.GetString(config.EnvPaymailAvatar)

	_, err := a.api.CreatePaymail(context.Background(), &commands.CreatePaymail{
		Key:        xpub,
		Address:    address,
		PublicName: alias,
		Avatar:     avatar,
	})
	if err != nil {
		a.log.Error().Msgf("Error while creating new paymail: %v", err.Error())
		return "", errors.Wrap(err, "error while creating new paymail")
	}

	return address, nil
}

func (a *adminClientAdapter) GetSharedConfig() (*models.SharedConfig, error) {
	sharedConfig, err := a.api.SharedConfig(context.Background())
	if err != nil {
		a.log.Error().Msgf("Error while fetching shared config: %v", err.Error())
		return nil, errors.Wrap(err, "error while fetching shared config")
	}

	return &models.SharedConfig{
		PaymailDomains:       sharedConfig.PaymailDomains,
		ExperimentalFeatures: sharedConfig.ExperimentalFeatures,
	}, nil
}

func newAdminClientAdapter(log *zerolog.Logger) (*adminClientAdapter, error) {
	adminKey := viper.GetString(config.EnvAdminXpriv)
	serverURL := viper.GetString(config.EnvServerURL)
	api, err := walletclient.NewAdminAPIWithXPriv(walletclientCfg.New(walletclientCfg.WithAddr(serverURL)), adminKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize admin API")
	}

	return &adminClientAdapter{api: api, log: log}, nil
}
