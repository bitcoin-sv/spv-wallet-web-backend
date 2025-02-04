package spvwallet

import (
	"context"
	"fmt"

	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	walletclientCfg "github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/common"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type userClientAdapter struct {
	api *walletclient.UserAPI
	log *zerolog.Logger
}

func (u *userClientAdapter) CreateAccessKey() (users.AccKey, error) {
	accessKey, err := u.api.GenerateAccessKey(context.Background(), &commands.GenerateAccessKey{})
	if err != nil {
		u.log.Error().Msgf("Error while creating new accessKey: %v", err.Error())
		return nil, errors.Wrap(err, "error while creating new accessKey ")
	}

	return &AccessKey{ID: accessKey.ID, Key: accessKey.Key}, nil
}

func (u *userClientAdapter) GetAccessKey(accessKeyID string) (users.AccKey, error) {
	accessKey, err := u.api.AccessKey(context.Background(), accessKeyID)
	if err != nil {
		u.log.Error().Str("accessKeyID", accessKeyID).Msgf("Error while getting accessKey: %v", err.Error())
		return nil, errors.Wrap(err, "error while getting accessKey")
	}

	return &AccessKey{ID: accessKey.ID, Key: accessKey.Key}, nil
}

func (u *userClientAdapter) RevokeAccessKey(accessKeyID string) (users.AccKey, error) {
	accessKey, err := u.api.AccessKey(context.Background(), accessKeyID)
	if err != nil {
		u.log.Error().Str("accessKeyID", accessKeyID).Msgf("Error while fetching accessKey: %v", err.Error())
		return nil, errors.Wrap(err, "error while fetching accessKey")
	}

	err = u.api.RevokeAccessKey(context.Background(), accessKeyID)
	if err != nil {
		u.log.Error().Str("accessKeyID", accessKeyID).Msgf("Error while revoking accessKey: %v", err.Error())
		return nil, errors.Wrap(err, "error while revoking accessKey")
	}

	return &AccessKey{ID: accessKey.ID, Key: accessKey.Key}, nil
}

// XPub Key methods
func (u *userClientAdapter) GetXPub() (users.PubKey, error) {
	xpub, err := u.api.XPub(context.Background())
	if err != nil {
		u.log.Error().Msgf("Error while getting new xPub: %v", err.Error())
		return nil, errors.Wrap(err, "error while getting new xPub")
	}

	return &XPub{ID: xpub.ID, CurrentBalance: xpub.CurrentBalance}, nil
}

func (u *userClientAdapter) SendToRecipients(recipients []*commands.Recipients, senderPaymail string) (users.Transaction, error) {
	// Send transaction.
	transaction, err := u.api.SendToRecipients(context.Background(), &commands.SendToRecipients{
		Recipients: recipients,
		Metadata: map[string]any{
			"receiver": recipients[0].To,
			"sender":   senderPaymail,
		},
	})
	if err != nil {
		u.log.Error().Msgf("Error while creating new tx: %v", err.Error())
		return nil, errors.Wrap(err, "error while creating new tx")
	}

	return &Transaction{
		ID:         transaction.ID,
		Direction:  fmt.Sprint(transaction.TransactionDirection),
		TotalValue: transaction.TotalValue,
		Status:     transaction.Status,
		CreatedAt:  transaction.Model.CreatedAt,
	}, nil
}

func (u *userClientAdapter) GetTransactions(queryParam *filter.QueryParams, userPaymail string) ([]users.Transaction, error) {
	if queryParam.OrderByField == "" {
		queryParam.OrderByField = "created_at"
	}

	if queryParam.SortDirection == "" {
		queryParam.SortDirection = "desc"
	}

	page, err := u.api.Transactions(context.Background(), queries.QueryWithPageFilter[filter.TransactionFilter](filter.Page{
		Number: queryParam.Page,
		Size:   queryParam.PageSize,
		Sort:   queryParam.SortDirection,
		SortBy: queryParam.OrderByField,
	}))
	if err != nil {
		u.log.Error().Str("userPaymail", userPaymail).Msgf("Error while getting transactions: %v", err.Error())
		return nil, errors.Wrap(err, "error while getting transactions")
	}

	var transactionsData = make([]users.Transaction, 0)
	for _, transaction := range page.Content {
		sender, receiver := GetPaymailsFromMetadata(transaction, userPaymail)
		status := "unconfirmed"
		if transaction.BlockHeight > 0 {
			status = "confirmed"
		}

		transactionsData = append(transactionsData, &Transaction{
			ID:         transaction.ID,
			Direction:  fmt.Sprint(transaction.TransactionDirection),
			TotalValue: getAbsoluteValue(transaction.OutputValue),
			Fee:        transaction.Fee,
			Status:     status,
			CreatedAt:  transaction.Model.CreatedAt,
			Sender:     sender,
			Receiver:   receiver,
		})
	}

	return transactionsData, nil
}

func (u *userClientAdapter) GetTransaction(transactionID, userPaymail string) (users.FullTransaction, error) {
	transaction, err := u.api.Transaction(context.Background(), transactionID)
	if err != nil {
		u.log.Error().Str("transactionId", transactionID).Str("userPaymail", userPaymail).Msgf("Error while getting transaction: %v", err.Error())
		return nil, errors.Wrap(err, "error while getting transaction")
	}

	sender, receiver := GetPaymailsFromMetadata(transaction, userPaymail)
	return &FullTransaction{
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
	}, nil
}

func (u *userClientAdapter) GetTransactionsCount() (int64, error) {
	return 0, nil // Note: Functionality it's not a part of the SPV Wallet Go client.
}

func (u *userClientAdapter) CreateAndFinalizeTransaction(recipients []*commands.Recipients, metadata map[string]any) (users.DraftTransaction, error) {
	draftTx, err := u.api.SendToRecipients(context.Background(), &commands.SendToRecipients{
		Recipients: recipients,
		Metadata:   metadata,
	})
	if err != nil {
		u.log.Error().Msgf("Error while sending to recipients: %v", err.Error())
		return nil, errors.Wrap(err, "error while sending to recipients")
	}

	return &DraftTransaction{
		TxDraftID: draftTx.DraftID,
		TxHex:     draftTx.Hex,
	}, nil
}

func (u *userClientAdapter) RecordTransaction(hex, draftTxID string, metadata map[string]any) (*models.Transaction, error) {
	tx, err := u.api.RecordTransaction(context.Background(), &commands.RecordTransaction{
		Metadata:    metadata,
		Hex:         hex,
		ReferenceID: draftTxID,
	})
	if err != nil {
		u.log.Error().Str("draftTxID", draftTxID).Msgf("Error while recording tx: %v", err.Error())
		return nil, errors.Wrap(err, "error while recording tx")
	}

	return &models.Transaction{
		Model:                common.Model(tx.Model),
		ID:                   tx.ID,
		Hex:                  tx.Hex,
		XpubInIDs:            tx.XpubInIDs,
		XpubOutIDs:           tx.XpubOutIDs,
		BlockHash:            tx.BlockHash,
		BlockHeight:          tx.BlockHeight,
		Fee:                  tx.Fee,
		NumberOfInputs:       tx.NumberOfInputs,
		NumberOfOutputs:      tx.NumberOfOutputs,
		DraftID:              tx.DraftID,
		TotalValue:           tx.TotalValue,
		OutputValue:          tx.OutputValue,
		Outputs:              tx.Outputs,
		Status:               tx.Status,
		TransactionDirection: tx.TransactionDirection,
	}, nil
}

// Contacts methods
func (u *userClientAdapter) UpsertContact(ctx context.Context, paymail, fullName, requesterPaymail string, metadata map[string]any) (*models.Contact, error) {
	contact, err := u.api.UpsertContact(ctx, commands.UpsertContact{
		ContactPaymail:   paymail,
		FullName:         fullName,
		Metadata:         metadata,
		RequesterPaymail: requesterPaymail,
	})
	if err != nil {
		return nil, errors.Wrap(err, "upsert contact error")
	}

	return &models.Contact{
		Model:    common.Model(contact.Model),
		ID:       contact.ID,
		FullName: contact.FullName,
		Paymail:  contact.Paymail,
		PubKey:   contact.PubKey,
		Status:   contact.Status,
	}, nil
}

func (u *userClientAdapter) AcceptContact(ctx context.Context, paymail string) error {
	return errors.Wrap(u.api.AcceptInvitation(ctx, paymail), "accept contact error")
}

func (u *userClientAdapter) RejectContact(ctx context.Context, paymail string) error {
	return errors.Wrap(u.api.RejectInvitation(ctx, paymail), "reject contact error")
}

func (u *userClientAdapter) ConfirmContact(ctx context.Context, contact *models.Contact, passcode, requesterPaymail string, period, digits uint) error {
	return errors.Wrap(u.api.ConfirmContact(ctx, contact, passcode, requesterPaymail, period, digits), "confirm contact error")
}

func (u *userClientAdapter) GetContacts(ctx context.Context, conditions *filter.ContactFilter, metadata map[string]any, queryParams *filter.QueryParams) (*models.SearchContactsResponse, error) {
	opts := []queries.QueryOption[filter.ContactFilter]{
		queries.QueryWithMetadataFilter[filter.ContactFilter](metadata),
	}

	if queryParams != nil {
		opts = append(opts,
			queries.QueryWithPageFilter[filter.ContactFilter](filter.Page{
				Number: queryParams.Page,
				Size:   queryParams.PageSize,
				Sort:   queryParams.SortDirection,
				SortBy: queryParams.OrderByField,
			}))
	}

	if conditions != nil {
		opts = append(opts,
			queries.QueryWithFilter(filter.ContactFilter{
				ModelFilter: conditions.ModelFilter,
				ID:          conditions.ID,
				FullName:    conditions.FullName,
				Paymail:     conditions.Paymail,
				PubKey:      conditions.PubKey,
				Status:      conditions.Status,
			}))
	}

	res, err := u.api.Contacts(ctx, opts...)
	if err != nil {
		u.log.Error().Msgf("Error while fetching contacts: %v", err.Error())
		return nil, errors.Wrap(err, "error while fetching contacts")
	}

	content := make([]*models.Contact, len(res.Content))
	for i, c := range res.Content {
		content[i] = &models.Contact{
			Model:    common.Model(c.Model),
			FullName: c.FullName,
			ID:       c.ID,
			Paymail:  c.Paymail,
			PubKey:   c.PubKey,
			Status:   c.Status,
		}
	}

	page := models.Page{
		TotalElements: int64(res.Page.TotalElements),
		TotalPages:    res.Page.TotalPages,
		Size:          res.Page.Size,
		Number:        res.Page.Number,
	}

	if queryParams != nil {
		page.OrderByField = &queryParams.OrderByField
		page.SortDirection = &queryParams.SortDirection
	}

	return &models.SearchContactsResponse{
		Content: content,
		Page:    page}, nil
}

func (u *userClientAdapter) GenerateTotpForContact(contact *models.Contact, period, digits uint) (string, error) {
	totp, err := u.api.GenerateTotpForContact(contact, period, digits)
	return totp, errors.Wrap(err, "error while generating TOTP for contact")
}

func newUserClientAdapterWithXPriv(log *zerolog.Logger, xPriv string) (*userClientAdapter, error) {
	serverURL := viper.GetString(config.EnvServerURL)
	api, err := walletclient.NewUserAPIWithXPriv(walletclientCfg.New(walletclientCfg.WithAddr(serverURL)), xPriv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize user API")
	}

	return &userClientAdapter{api: api, log: log}, nil
}

func newUserClientAdapterWithAccessKey(log *zerolog.Logger, accessKey string) (*userClientAdapter, error) {
	serverURL := viper.GetString(config.EnvServerURL)
	api, err := walletclient.NewUserAPIWithAccessKey(walletclientCfg.New(walletclientCfg.WithAddr(serverURL)), accessKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize user API")
	}

	return &userClientAdapter{api: api, log: log}, nil
}
