package buxclient

import "github.com/libsv/go-bk/bip32"

type (
	// UserBuxClient defines methods for bux client with user key.
	UserBuxClient interface {
		// Access Key methods
		CreateAccessKey() (*AccessKey, error)
		GetAccessKey(accessKeyId string) (*AccessKey, error)
		RevokeAccessKey(accessKeyId string) (*AccessKey, error)
	}

	// AdmBuxClient defines methods for bux client with admin key.
	AdmBuxClient interface {
		RegisterXpub(xpriv *bip32.ExtendedKey) (string, error)
		RegisterPaymail(alias, xpub string) (string, error)
	}

	// BuxClientFactory defines methods for bux client factory.
	BuxClientFactory interface {
		CreateWithXpriv(xpriv string) (UserBuxClient, error)
		CreateAdminBuxClient() (AdmBuxClient, error)
	}
)
