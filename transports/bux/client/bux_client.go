package buxclient

import (
	"context"
	"fmt"
	"math"

	"bux-wallet/domain/users"
	"bux-wallet/logging"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/go-buxclient"
	"github.com/BuxOrg/go-buxclient/transports"
	"github.com/mrz1836/go-datastore"
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

// SendToRecipients sends satoshis to recipients.
func (c *BuxClient) SendToRecipients(recipients []*transports.Recipients, senderPaymail string) (users.Transaction, error) {
	// Create matadata with sender and receiver paymails.
	metadata := &bux.Metadata{
		"receiver": recipients[0].To,
		"sender":   senderPaymail,
	}

	// Send transaction.
	transaction, err := c.client.SendToRecipients(context.Background(), recipients, metadata)
	if err != nil {
		return nil, err
	}

	t := &Transaction{
		Id:         transaction.ID,
		Direction:  fmt.Sprint(transaction.Direction),
		TotalValue: transaction.TotalValue,
		Status:     transaction.Status.String(),
		CreatedAt:  transaction.CreatedAt,
	}
	return t, nil
}

// CreateAndFinalizeTransaction creates draft transaction and finalizes it.
func (c *BuxClient) CreateAndFinalizeTransaction(recipients []*transports.Recipients, metadata *bux.Metadata) (users.DraftTransaction, error) {
	// Create draft transaction.
	draftTx, err := c.client.DraftToRecipients(context.Background(), recipients, metadata)
	if err != nil {
		return nil, err
	}

	// Finalize draft transaction.
	hex, err := c.client.FinalizeTransaction(draftTx)
	if err != nil {
		return nil, err
	}

	draftTransaction := DraftTransaction{
		TxDraftId: draftTx.ID,
		TxHex:     hex,
	}

	return &draftTransaction, nil
}

// RecordTransaction records transaction in BUX.
func (c *BuxClient) RecordTransaction(hex, draftTxId string, metadata *bux.Metadata) {
	c.client.RecordTransaction(context.Background(), hex, draftTxId, metadata) // nolint: all // TODO: handle error in correct way.
}

// GetTransactions returns all transactions.
func (c *BuxClient) GetTransactions(queryParam datastore.QueryParams, userPaymail string) ([]users.Transaction, error) {
	conditions := make(map[string]interface{})

	if queryParam.OrderByField == "" {
		queryParam.OrderByField = "created_at"
	}

	if queryParam.SortDirection == "" {
		queryParam.SortDirection = "desc"
	}

	transactions, err := c.client.GetTransactions(context.Background(), conditions, &bux.Metadata{}, &queryParam)
	if err != nil {
		return nil, err
	}

	var transactionsData = make([]users.Transaction, 0)
	for _, transaction := range transactions {
		sender, receiver := getPaymailsFromMetadata(transaction, userPaymail)
		absTotalValue := uint64(math.Abs(float64(transaction.OutputValue)))
		status := "unconfirmed"
		if transaction.BlockHeight > 0 {
			status = "confirmed"
		}
		transactionData := Transaction{
			Id:         transaction.ID,
			Direction:  fmt.Sprint(transaction.Direction),
			TotalValue: absTotalValue,
			Fee:        transaction.Fee,
			Status:     status,
			CreatedAt:  transaction.CreatedAt,
			Sender:     sender,
			Receiver:   receiver,
		}
		transactionsData = append(transactionsData, &transactionData)
	}

	return transactionsData, nil
}

// GetTransaction returns transaction by id.
func (c *BuxClient) GetTransaction(transactionId, userPaymail string) (users.FullTransaction, error) {
	transaction, err := c.client.GetTransaction(context.Background(), transactionId)
	if err != nil {
		return nil, err
	}

	sender, receiver := getPaymailsFromMetadata(transaction, userPaymail)
	absTotalValue := uint64(math.Abs(float64(transaction.OutputValue)))

	transactionData := FullTransaction{
		Id:              transaction.ID,
		BlockHash:       transaction.BlockHash,
		BlockHeight:     transaction.BlockHeight,
		TotalValue:      absTotalValue,
		Direction:       fmt.Sprint(transaction.Direction),
		Status:          transaction.Status.String(),
		Fee:             transaction.Fee,
		NumberOfInputs:  transaction.NumberOfInputs,
		NumberOfOutputs: transaction.NumberOfOutputs,
		CreatedAt:       transaction.CreatedAt,
		Sender:          sender,
		Receiver:        receiver,
	}

	return &transactionData, nil
}

// GetTransactionsCount returns number of transactions.
func (c *BuxClient) GetTransactionsCount() (int64, error) {
	conditions := make(map[string]interface{})

	count, err := c.client.GetTransactionsCount(context.Background(), conditions, &bux.Metadata{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// getPaymailsFromMetadata returns sender and receiver paymails from metadata.
// If no paymail was found in metadata, fallback paymail is returned.
func getPaymailsFromMetadata(transaction *bux.Transaction, fallbackPaymail string) (string, string) {
	senderPaymail := ""
	receiverPaymail := ""

	if transaction.Metadata != nil {
		// Try to get paymails from metadata if the transaction was made in BUX.
		if transaction.Metadata["sender"] != nil {
			senderPaymail = transaction.Metadata["sender"].(string)
		}
		if transaction.Metadata["receiver"] != nil {
			receiverPaymail = transaction.Metadata["receiver"].(string)
		}

		if senderPaymail == "" {
			// Try to get paymails from metadata if the transaction was made outside BUX.
			if transaction.Metadata["p2p_tx_metadata"] != nil {
				p2pTxMetadata := transaction.Metadata["p2p_tx_metadata"].(map[string]interface{})
				if p2pTxMetadata["sender"] != nil {
					senderPaymail = p2pTxMetadata["sender"].(string)
				}
			}
		}
	}

	if transaction.Direction == "incoming" && receiverPaymail == "" {
		receiverPaymail = fallbackPaymail
	} else if transaction.Direction == "outgoing" && senderPaymail == "" {
		senderPaymail = fallbackPaymail
	}

	return senderPaymail, receiverPaymail
}
