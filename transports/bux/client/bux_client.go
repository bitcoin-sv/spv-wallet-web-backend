package buxclient

import (
	"context"

	"bux-wallet/domain/users"
	"bux-wallet/logging"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/go-buxclient"
)

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

// GetAccessKey returns access key.
func (a *AccessKey) GetAccessKey() string {
	return a.Key
}

// GetAccessKeyId returns access key id.
func (a *AccessKey) GetAccessKeyId() string {
	return a.Id
}

// CreateAccessKey creates new access key for user.
func (c *BuxClient) CreateAccessKey() (users.AccKey, error) {
	accessKey, err := c.client.CreateAccessKey(context.Background(), &bux.Metadata{})
	if err != nil {
		return nil, err
	}

	accessKeyData := AccessKey{
		Id:  accessKey.ID,
		Key: accessKey.Key,
	}

	return &accessKeyData, err
}

// GetAccessKey checks if access key is valid.
func (c *BuxClient) GetAccessKey(accessKeyId string) (users.AccKey, error) {
	accessKey, err := c.client.GetAccessKey(context.Background(), accessKeyId)
	if err != nil {
		return nil, err
	}

	accessKeyData := AccessKey{
		Id:  accessKey.ID,
		Key: accessKey.Key,
	}

	return &accessKeyData, nil
}

// RevokeAccessKey revokes access key.
func (c *BuxClient) RevokeAccessKey(accessKeyId string) (users.AccKey, error) {
	accessKey, err := c.client.RevokeAccessKey(context.Background(), accessKeyId)
	if err != nil {
		return nil, err
	}

	accessKeyData := AccessKey{
		Id:  accessKey.ID,
		Key: accessKey.Key,
	}

	return &accessKeyData, nil
}
