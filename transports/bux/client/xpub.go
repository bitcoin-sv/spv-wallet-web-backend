package buxclient

import (
	"context"
	"fmt"

	"github.com/BuxOrg/bux"
	"github.com/libsv/go-bk/bip32"
)

// RegisterXpub registers xpub in bux.
func (adminBuxClient *AdminBuxClient) RegisterXpub(xpriv *bip32.ExtendedKey) error {
	// Get xpub from xpriv.
	xpub, err := xpriv.Neuter()

	if err != nil {
		fmt.Println(err)
		return err
	}

	// Register new xpub in BUX.
	err = adminBuxClient.client.NewXpub(
		context.Background(), xpub.String(), &bux.Metadata{},
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
