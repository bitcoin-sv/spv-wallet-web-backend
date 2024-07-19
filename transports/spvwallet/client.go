package spvwallet

import (
	"context"
	"fmt"
	"math"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Client implements UserWalletClient interface which wraps the spv-wallet-go-client and provides methods for user.
type Client struct {
	client *walletclient.WalletClient
	log    *zerolog.Logger
}

// CreateAccessKey creates new access key for user.
func (c *Client) CreateAccessKey() (users.AccKey, error) {
	accessKey, err := c.client.CreateAccessKey(context.Background(), nil)
	if err != nil {
		c.log.Error().Msgf("Error while creating new accessKey: %v", err.Error())
		return nil, err
	}

	accessKeyData := AccessKey{
		ID:  accessKey.ID,
		Key: accessKey.Key,
	}

	return &accessKeyData, err
}

// GetAccessKey checks if access key is valid.
func (c *Client) GetAccessKey(accessKeyID string) (users.AccKey, error) {
	accessKey, err := c.client.GetAccessKey(context.Background(), accessKeyID)
	if err != nil {
		c.log.Error().
			Str("accessKeyID", accessKeyID).
			Msgf("Error while getting accessKey: %v", err.Error())
		return nil, err
	}

	accessKeyData := AccessKey{
		ID:  accessKey.ID,
		Key: accessKey.Key,
	}

	return &accessKeyData, nil
}

// RevokeAccessKey revokes access key.
func (c *Client) RevokeAccessKey(accessKeyID string) (users.AccKey, error) {
	accessKey, err := c.client.RevokeAccessKey(context.Background(), accessKeyID)
	if err != nil {
		c.log.Error().
			Str("accessKeyID", accessKeyID).
			Msgf("Error while revoking accessKey: %v", err.Error())
		return nil, err
	}

	accessKeyData := AccessKey{
		ID:  accessKey.ID,
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
		ID:             xpub.ID,
		CurrentBalance: xpub.CurrentBalance,
	}

	return &xPub, nil
}

// SendToRecipients sends satoshis to recipients.
func (c *Client) SendToRecipients(recipients []*walletclient.Recipients, senderPaymail string) (users.Transaction, error) {
	// Create matadata with sender and receiver paymails.
	metadata := map[string]any{
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
		ID:         transaction.ID,
		Direction:  fmt.Sprint(transaction.TransactionDirection),
		TotalValue: transaction.TotalValue,
		Status:     transaction.Status,
		CreatedAt:  transaction.Model.CreatedAt,
	}
	return t, nil
}

// CreateAndFinalizeTransaction creates draft transaction and finalizes it.
func (c *Client) CreateAndFinalizeTransaction(recipients []*walletclient.Recipients, metadata map[string]any) (users.DraftTransaction, error) {
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
		TxDraftID: draftTx.ID,
		TxHex:     hex,
	}

	return &draftTransaction, nil
}

// RecordTransaction records transaction in SPV Wallet.
func (c *Client) RecordTransaction(hex, draftTxID string, metadata map[string]any) (*models.Transaction, error) {
	tx, err := c.client.RecordTransaction(context.Background(), hex, draftTxID, metadata)
	if err != nil {
		c.log.Error().Str("draftTxID", draftTxID).Msgf("Error while recording tx: %v", err.Error())
		return nil, err
	}
	return tx, nil
}

// GetTransactions returns all transactions.
func (c *Client) GetTransactions(queryParam *filter.QueryParams, userPaymail string) ([]users.Transaction, error) {
	if queryParam.OrderByField == "" {
		queryParam.OrderByField = "created_at"
	}

	if queryParam.SortDirection == "" {
		queryParam.SortDirection = "desc"
	}

	transactions, err := c.client.GetTransactions(context.Background(), nil, nil, queryParam)
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
			ID:         transaction.ID,
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
func (c *Client) GetTransaction(transactionID, userPaymail string) (users.FullTransaction, error) {
	transaction, err := c.client.GetTransaction(context.Background(), transactionID)
	if err != nil {
		c.log.Error().
			Str("transactionId", transactionID).
			Str("userPaymail", userPaymail).
			Msgf("Error while getting transaction: %v", err.Error())
		return nil, err
	}

	sender, receiver := GetPaymailsFromMetadata(transaction, userPaymail)

	transactionData := FullTransaction{
		ID:              transaction.ID,
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
	count, err := c.client.GetTransactionsCount(context.Background(), nil, nil)
	if err != nil {
		c.log.Error().Msgf("Error while getting transactions count: %v", err.Error())
		return 0, err
	}
	return count, nil
}

// UpsertContact creates or updates contact.
func (c *Client) UpsertContact(ctx context.Context, paymail, fullName, requesterPaymail string, metadata map[string]any) (*models.Contact, error) {
	return c.client.UpsertContact(ctx, paymail, fullName, requesterPaymail, metadata)
}

// AcceptContact accepts contact.
func (c *Client) AcceptContact(ctx context.Context, paymail string) error {
	return c.client.AcceptContact(ctx, paymail)
}

// RejectContact rejects contact.
func (c *Client) RejectContact(ctx context.Context, paymail string) error {
	return c.client.RejectContact(ctx, paymail)
}

// ConfirmContact confirms contact.
func (c *Client) ConfirmContact(ctx context.Context, contact *models.Contact, passcode, requesterPaymail string, period, digits uint) error {
	return c.client.ConfirmContact(ctx, contact, passcode, requesterPaymail, period, digits)
}

// GetContacts returns all contacts.
func (c *Client) GetContacts(ctx context.Context, conditions *filter.ContactFilter, metadata map[string]any, queryParams *filter.QueryParams) (*models.SearchContactsResponse, error) {
	return c.client.GetContacts(ctx, conditions, metadata, queryParams)
}

// GenerateTotpForContact generates TOTP for contact.
func (c *Client) GenerateTotpForContact(contact *models.Contact, period, digits uint) (string, error) {
	totp, err := c.client.GenerateTotpForContact(contact, period, digits)
	return totp, errors.Wrap(err, "error while generating TOTP for contact")
}

func getAbsoluteValue(value int64) uint64 {
	return uint64(math.Abs(float64(value)))
}
