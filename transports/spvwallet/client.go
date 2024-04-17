package spvwallet

import (
	"context"
	"fmt"
	"math"

	"github.com/rs/zerolog"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
	"github.com/bitcoin-sv/spv-wallet/models"

	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
)

// Client implements UserWalletClient interface which wraps the spv-wallet-go-client and provides methods for user.
type Client struct {
	client *walletclient.WalletClient
	log    *zerolog.Logger
}

// CreateAccessKey creates new access key for user.
func (c *Client) CreateAccessKey() (users.AccKey, error) {
	accessKey, err := c.client.CreateAccessKey(context.Background(), &models.Metadata{})
	if err != nil {
		c.log.Error().Msgf("Error while creating new accessKey: %v", err.Error())
		return nil, err
	}

	accessKeyData := AccessKey{
		Id:  accessKey.ID,
		Key: accessKey.Key,
	}

	return &accessKeyData, err
}

// GetAccessKey checks if access key is valid.
func (c *Client) GetAccessKey(accessKeyId string) (users.AccKey, error) {
	accessKey, err := c.client.GetAccessKey(context.Background(), accessKeyId)
	if err != nil {
		c.log.Error().
			Str("accessKeyID", accessKeyId).
			Msgf("Error while getting accessKey: %v", err.Error())
		return nil, err
	}

	accessKeyData := AccessKey{
		Id:  accessKey.ID,
		Key: accessKey.Key,
	}

	return &accessKeyData, nil
}

// RevokeAccessKey revokes access key.
func (c *Client) RevokeAccessKey(accessKeyId string) (users.AccKey, error) {
	accessKey, err := c.client.RevokeAccessKey(context.Background(), accessKeyId)
	if err != nil {
		c.log.Error().
			Str("accessKeyID", accessKeyId).
			Msgf("Error while revoking accessKey: %v", err.Error())
		return nil, err
	}

	accessKeyData := AccessKey{
		Id:  accessKey.ID,
		Key: accessKey.Key,
	}

	return &accessKeyData, nil
}

// GetXPub returns xpub.
func (c *Client) GetXPub() (users.PubKey, error) {
	xpub, err := c.client.GetXPub(context.Background())
	if err != nil {
		c.log.Error().Msgf("Error while getting new xPub: %v", err.Error())
		return nil, err
	}

	xPub := XPub{
		Id:             xpub.ID,
		CurrentBalance: xpub.CurrentBalance,
	}

	return &xPub, nil
}

// SendToRecipients sends satoshis to recipients.
func (c *Client) SendToRecipients(recipients []*transports.Recipients, senderPaymail string) (users.Transaction, error) {
	// Create matadata with sender and receiver paymails.
	metadata := &models.Metadata{
		"receiver": recipients[0].To,
		"sender":   senderPaymail,
	}

	// Send transaction.
	transaction, err := c.client.SendToRecipients(context.Background(), recipients, metadata)
	if err != nil {
		c.log.Error().Msgf("Error while creating new tx: %v", err.Error())
		return nil, err
	}

	t := &Transaction{
		Id:         transaction.ID,
		Direction:  fmt.Sprint(transaction.TransactionDirection),
		TotalValue: transaction.TotalValue,
		Status:     transaction.Status,
		CreatedAt:  transaction.Model.CreatedAt,
	}
	return t, nil
}

// CreateAndFinalizeTransaction creates draft transaction and finalizes it.
func (c *Client) CreateAndFinalizeTransaction(recipients []*transports.Recipients, metadata *models.Metadata) (users.DraftTransaction, error) {
	// Create draft transaction.
	draftTx, err := c.client.DraftToRecipients(context.Background(), recipients, metadata)
	if err != nil {
		c.log.Error().Msgf("Error while creating new draft tx: %v", err.Error())
		return nil, err
	}

	// Finalize draft transaction.
	hex, err := c.client.FinalizeTransaction(draftTx)
	if err != nil {
		c.log.Error().Str("draftTxID", draftTx.ID).Msgf("Error while finalizing tx: %v", err.Error())
		return nil, err
	}

	draftTransaction := DraftTransaction{
		TxDraftId: draftTx.ID,
		TxHex:     hex,
	}

	return &draftTransaction, nil
}

// RecordTransaction records transaction in SPV Wallet.
func (c *Client) RecordTransaction(hex, draftTxId string, metadata *models.Metadata) (*models.Transaction, error) {
	tx, err := c.client.RecordTransaction(context.Background(), hex, draftTxId, metadata)
	if err != nil {
		c.log.Error().Str("draftTxID", draftTxId).Msgf("Error while recording tx: %v", err.Error())
		return nil, err
	}
	return tx, nil
}

// GetTransactions returns all transactions.
func (c *Client) GetTransactions(queryParam transports.QueryParams, userPaymail string) ([]users.Transaction, error) {
	conditions := make(map[string]interface{})

	if queryParam.OrderByField == "" {
		queryParam.OrderByField = "created_at"
	}

	if queryParam.SortDirection == "" {
		queryParam.SortDirection = "desc"
	}

	transactions, err := c.client.GetTransactions(context.Background(), conditions, &models.Metadata{}, &queryParam)
	if err != nil {
		c.log.Error().
			Str("userPaymail", userPaymail).
			Msgf("Error while getting transactions: %v", err.Error())
		return nil, err
	}

	var transactionsData = make([]users.Transaction, 0)
	for _, transaction := range transactions {
		sender, receiver := GetPaymailsFromMetadata(transaction, userPaymail)
		status := "unconfirmed"
		if transaction.BlockHeight > 0 {
			status = "confirmed"
		}
		transactionData := Transaction{
			Id:         transaction.ID,
			Direction:  fmt.Sprint(transaction.TransactionDirection),
			TotalValue: getAbsoluteValue(transaction.OutputValue),
			Fee:        transaction.Fee,
			Status:     status,
			CreatedAt:  transaction.Model.CreatedAt,
			Sender:     sender,
			Receiver:   receiver,
		}
		transactionsData = append(transactionsData, &transactionData)
	}

	return transactionsData, nil
}

// GetTransaction returns transaction by id.
func (c *Client) GetTransaction(transactionId, userPaymail string) (users.FullTransaction, error) {
	transaction, err := c.client.GetTransaction(context.Background(), transactionId)
	if err != nil {
		c.log.Error().
			Str("transactionId", transactionId).
			Str("userPaymail", userPaymail).
			Msgf("Error while getting transaction: %v", err.Error())
		return nil, err
	}

	sender, receiver := GetPaymailsFromMetadata(transaction, userPaymail)

	transactionData := FullTransaction{
		Id:              transaction.ID,
		BlockHash:       transaction.BlockHash,
		BlockHeight:     transaction.BlockHeight,
		TotalValue:      getAbsoluteValue(transaction.OutputValue),
		Direction:       fmt.Sprint(transaction.TransactionDirection),
		Status:          transaction.Status,
		Fee:             transaction.Fee,
		NumberOfInputs:  transaction.NumberOfInputs,
		NumberOfOutputs: transaction.NumberOfOutputs,
		CreatedAt:       transaction.Model.CreatedAt,
		Sender:          sender,
		Receiver:        receiver,
	}

	return &transactionData, nil
}

// GetTransactionsCount returns number of transactions.
func (c *Client) GetTransactionsCount() (int64, error) {
	conditions := make(map[string]interface{})

	count, err := c.client.GetTransactionsCount(context.Background(), conditions, &models.Metadata{})
	if err != nil {
		c.log.Error().Msgf("Error while getting transactions count: %v", err.Error())
		return 0, err
	}
	return count, nil
}

// UpsertContact creates or updates contact.
func (c *Client) UpsertContact(ctx context.Context, paymail, fullName string, metadata *models.Metadata) (*models.Contact, transports.ResponseError) {
	return c.client.UpsertContact(ctx, paymail, fullName, metadata)
}

// AcceptContact accepts contact.
func (c *Client) AcceptContact(ctx context.Context, paymail string) transports.ResponseError {
	return c.client.AcceptContact(ctx, paymail)
}

// RejectContact rejects contact.
func (c *Client) RejectContact(ctx context.Context, paymail string) transports.ResponseError {
	return c.client.RejectContact(ctx, paymail)
}

// ConfirmContact confirms contact.
func (c *Client) ConfirmContact(ctx context.Context, contact *models.Contact, passcode string, period, digits uint) transports.ResponseError {
	return c.client.ConfirmContact(ctx, contact, passcode, period, digits)
}

// GetContacts returns all contacts.
func (c *Client) GetContacts(ctx context.Context, conditions map[string]interface{}, metadata *models.Metadata, queryParams *transports.QueryParams) ([]*models.Contact, transports.ResponseError) {
	return c.client.GetContacts(ctx, conditions, metadata, queryParams)
}

// GenerateTotpForContact generates TOTP for contact.
func (c *Client) GenerateTotpForContact(contact *models.Contact, period, digits uint) (string, error) {
	return c.client.GenerateTotpForContact(contact, period, digits)
}

func getAbsoluteValue(value int64) uint64 {
	return uint64(math.Abs(float64(value)))
}
