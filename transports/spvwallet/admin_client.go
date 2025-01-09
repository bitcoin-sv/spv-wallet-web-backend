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
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type AdminClientAdapter struct {
	log *zerolog.Logger
	api *walletclient.AdminAPI
}

func (a *AdminClientAdapter) RegisterXpub(xpriv *bip32.ExtendedKey) (string, error) {
	// Get xpub from xpriv.
	xpub, err := xpriv.Neuter()
	if err != nil {
		a.log.Error().Msgf("Error while creating new xPub: %v", err.Error())
		return "", err
	}

	_, err = a.api.CreateXPub(context.TODO(), &commands.CreateUserXpub{XPub: xpriv.String()})
	if err != nil {
		a.log.Error().Str("xpub", xpub.String()).Msgf("Error while registering new xPub: %v", err.Error())
		return "", err
	}

	return xpub.String(), nil
}

func (a *AdminClientAdapter) RegisterPaymail(alias, xpub string) (string, error) {
	// Get paymail domain from env.
	domain := viper.GetString(config.EnvPaymailDomain)

	// Create paymail address.
	address := fmt.Sprintf("%s@%s", alias, domain)

	// Get avatar url from env.
	avatar := viper.GetString(config.EnvPaymailAvatar)

	_, err := a.api.CreatePaymail(context.TODO(), &commands.CreatePaymail{
		Key:        xpub,
		Address:    address,
		PublicName: alias,
		Avatar:     avatar,
	})
	if err != nil {
		a.log.Error().Msgf("Error while registering new paymail: %v", err.Error())
		return "", err
	}

	return address, nil
}

func (a *AdminClientAdapter) GetSharedConfig() (*models.SharedConfig, error) {
	sharedConfig, err := a.api.SharedConfig(context.TODO())
	if err != nil {
		a.log.Error().Msgf("Error while getting shared config: %v", err.Error())
		return nil, err
	}

	return &models.SharedConfig{
		PaymailDomains:       sharedConfig.PaymailDomains,
		ExperimentalFeatures: sharedConfig.ExperimentalFeatures,
	}, nil
}

func NewAdminClientAdapter(log *zerolog.Logger) (*AdminClientAdapter, error) {
	adminKey := viper.GetString(config.EnvAdminXpriv)
	serverURL := viper.GetString(config.EnvServerURL)
	api, err := walletclient.NewAdminAPIWithXPriv(walletclientCfg.New(walletclientCfg.WithAddr(serverURL)), adminKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize admin API: %w", err)
	}

	return &AdminClientAdapter{api: api, log: log}, nil
}
