package buxclient

import (
	"bux-wallet/config"
	"context"
	"fmt"

	"github.com/BuxOrg/bux"
	"github.com/spf13/viper"
)

// RegisterPaymail registers new paymail in bux.
func (adminBuxClient *AdminBuxClient) RegisterPaymail(alias, xpub string) (string, error) {
	// Get paymail domain from env.
	domain := viper.GetString(config.EnvBuxPaymailDomain)

	// Create paymail address.
	address := fmt.Sprintf("%s@%s", alias, domain)

	// Register new xpub in BUX.
	err := adminBuxClient.client.NewPaymail(context.Background(), xpub, address, alias, alias, &bux.Metadata{})

	if err != nil {
		return "", err
	}
	return address, nil
}
