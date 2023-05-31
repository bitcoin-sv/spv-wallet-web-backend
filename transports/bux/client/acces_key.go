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
