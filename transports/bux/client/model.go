package buxclient

import (
	"bux-wallet/logging"

	"github.com/BuxOrg/go-buxclient"
)

// AdminBuxClient is a wrapper for Admin Bux Client.
type AdminBuxClient struct {
	client *buxclient.BuxClient
	log    logging.Logger
}

// BuxClient is a wrapper for Bux Client.
type BuxClient struct {
	client *buxclient.BuxClient
	log    logging.Logger
}

// AccessKey is a struct that contains access key data.
type AccessKey struct {
	Id  string `json:"id"`
	Key string `json:"key"`
}
