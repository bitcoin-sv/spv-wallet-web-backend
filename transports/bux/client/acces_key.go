package buxclient

import (
	"context"

	"github.com/BuxOrg/bux"
)

// CreateAccessKey creates new access key for user.
func (c *BuxClient) CreateAccessKey() (*AccessKey, error) {
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
func (c *BuxClient) GetAccessKey(accessKeyId string) (*AccessKey, error) {
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
func (c *BuxClient) RevokeAccessKey(accessKeyId string) (*AccessKey, error) {
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
