package buxclient

import (
	"context"
	"fmt"

	"bux-wallet/domain/users"
	"bux-wallet/logging"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/go-buxclient"
	"github.com/BuxOrg/go-buxclient/transports"
)

// BuxClient is a wrapper for Bux Client.
type BuxClient struct {
	client *buxclient.BuxClient
	log    logging.Logger
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

// GetXPub returns xpub.
func (c *BuxClient) GetXPub() (users.PubKey, error) {
	xpub, err := c.client.GetXPub(context.Background())
	if err != nil {
		return nil, err
	}

	xPub := XPub{
		Id:             xpub.ID,
		XPub:           xpub.Model.RawXpub(),
		CurrentBalance: xpub.CurrentBalance,
	}

	return &xPub, nil
}

// SendToRecipents sends satoshis to recipients.
func (c *BuxClient) SendToRecipents(recipients []*transports.Recipients) (string, error) {
	transaction, err := c.client.SendToRecipients(context.Background(), recipients, &bux.Metadata{})
	if err != nil {
		return "", err
	}
	return transaction.ID, nil
}

// GetTransactions returns all transactions.
func (c *BuxClient) GetTransactions() ([]users.Transaction, error) {
	conditions := make(map[string]interface{})
	transactions, err := c.client.GetTransactions(context.Background(), conditions, &bux.Metadata{})
	if err != nil {
		return nil, err
	}

	var transactionsData = make([]users.Transaction, 0)
	for _, transaction := range transactions {
		transactionData := Transaction{
			Id:         transaction.ID,
			Direction:  fmt.Sprint(transaction.Direction),
			TotalValue: transaction.TotalValue,
		}
		transactionsData = append(transactionsData, &transactionData)
	}

	return transactionsData, nil
}

// GetTransaction returns transaction by id.
func (c *BuxClient) GetTransaction(transactionId string) (users.FullTransaction, error) {
	transaction, err := c.client.GetTransaction(context.Background(), transactionId)
	if err != nil {
		return nil, err
	}

	transactionData := FullTransaction{
		Id:              transaction.ID,
		BlockHash:       transaction.BlockHash,
		BlockHeight:     transaction.BlockHeight,
		TotalValue:      transaction.TotalValue,
		Direction:       fmt.Sprint(transaction.Direction),
		Status:          transaction.Status.String(),
		Fee:             transaction.Fee,
		NumberOfInputs:  transaction.NumberOfInputs,
		NumberOfOutputs: transaction.NumberOfOutputs,
		CreatedAt:       transaction.CreatedAt,
	}

	return &transactionData, nil
}
