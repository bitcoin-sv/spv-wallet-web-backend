package buxclient

import (
	"context"

	"github.com/BuxOrg/bux"
)

// CreateAccessKey creates new access key for user.
func (c *BuxClient) CreateAccessKey() (string, error) {
	accessKey, err := c.client.CreateAccessKey(context.Background(), &bux.Metadata{})
	if err != nil {
		return "", err
	}

	return accessKey.ID, err
}

// GetAccessKey checks if access key is valid.
func (c *BuxClient) GetAccessKey(accessKeyId string) error {
	_, err := c.client.GetAccessKey(context.Background(), accessKeyId)
	if err != nil {
		return err
	}

	return nil
}

// RevokeAccessKey revokes access key.
func (c *BuxClient) RevokeAccessKey(accessKeyId string) error {
	_, err := c.client.RevokeAccessKey(context.Background(), accessKeyId)
	if err != nil {
		return err
	}

	return nil
}
