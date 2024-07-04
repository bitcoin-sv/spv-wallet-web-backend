package transactions

import (
	"math"
	"time"

	"github.com/avast/retry-go/v4"
	walletclient "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet-web-backend/notification"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// TransactionService represents service whoch contains methods linked with transactions.
type TransactionService struct {
	adminWalletClient   users.AdminWalletClient
	walletClientFactory users.WalletClientFactory
	log                 *zerolog.Logger
}

// NewTransactionService creates new transaction service.
func NewTransactionService(adminWalletClient users.AdminWalletClient, walletClientFactory users.WalletClientFactory, log *zerolog.Logger) *TransactionService {
	transactionServiceLogger := log.With().Str("service", "transaction-service").Logger()
	return &TransactionService{
		adminWalletClient:   adminWalletClient,
		walletClientFactory: walletClientFactory,
		log:                 &transactionServiceLogger,
	}
}

// CreateTransaction creates transaction.
func (s *TransactionService) CreateTransaction(userPaymail, xpriv, recipient string, satoshis uint64, events chan notification.TransactionEvent) error {
	userWalletClient, err := s.walletClientFactory.CreateWithXpriv(xpriv)
	if err != nil {
		return errors.Wrap(err, "spv wallet error")
	}

	var recipients = []*walletclient.Recipients{
		{
			Satoshis: satoshis,
			To:       recipient,
		},
	}

	metadata := map[string]any{
		"receiver": recipient,
		"sender":   userPaymail,
	}

	draftTransaction, err := userWalletClient.CreateAndFinalizeTransaction(recipients, metadata)
	if err != nil {
		return errors.Wrap(err, "spv wallet error")
	}

	go func() {
		tx, err := tryRecordTransaction(userWalletClient, draftTransaction, metadata, s.log)
		if err != nil {
			events <- notification.PrepareTransactionErrorEvent(err)
		} else if tx != nil {
			events <- notification.PrepareTransactionEvent(tx)
		}
	}()

	return nil
}

// GetTransaction returns transaction by id.
func (s *TransactionService) GetTransaction(accessKey, id, userPaymail string) (users.FullTransaction, error) {
	// Try to generate user-client with decrypted xpriv.
	userWalletClient, err := s.walletClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return nil, errors.Wrap(err, "spv wallet error")
	}

	transaction, err := userWalletClient.GetTransaction(id, userPaymail)
	if err != nil {
		return nil, errors.Wrap(err, "spv wallet error")
	}

	return transaction, nil
}

// GetTransactions returns transactions by access key.
func (s *TransactionService) GetTransactions(accessKey, userPaymail string, queryParam *filter.QueryParams) (*PaginatedTransactions, error) {
	// Try to generate user-client with decrypted xpriv.
	userWalletClient, err := s.walletClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return nil, errors.Wrap(err, "spv wallet error")
	}

	count, err := userWalletClient.GetTransactionsCount()
	if err != nil {
		return nil, errors.Wrap(err, "spv wallet error")
	}

	transactions, err := userWalletClient.GetTransactions(queryParam, userPaymail)
	if err != nil {
		return nil, errors.Wrap(err, "spv wallet error")
	}

	// Calculate pages.
	pages := int(math.Ceil(float64(count) / float64(queryParam.PageSize)))

	pTransactions := &PaginatedTransactions{
		Count:        count,
		Pages:        pages,
		Transactions: transactions,
	}

	return pTransactions, nil
}

func tryRecordTransaction(userWalletClient users.UserWalletClient, draftTx users.DraftTransaction, metadata map[string]any, log *zerolog.Logger) (*models.Transaction, error) {
	retries := uint(3)
	tx, recordErr := tryRecord(userWalletClient, draftTx, metadata, log, retries)

	if recordErr != nil {
		log.Error().
			Str("draftTxID", draftTx.GetDraftTransactionID()).
			Msgf("record transaction failed: %s", recordErr.Error())
		return nil, recordErr
	}

	log.Debug().
		Str("draftTxID", draftTx.GetDraftTransactionID()).
		Msg("transaction successfully recorded")
	return tx, nil
}

func tryRecord(userWalletClient users.UserWalletClient, draftTx users.DraftTransaction, metadata map[string]any, log *zerolog.Logger, retries uint) (*models.Transaction, error) {
	log.Debug().
		Str("draftTxID", draftTx.GetDraftTransactionID()).
		Msg("record transaction")

	tx := &models.Transaction{}
	err := retry.Do(
		func() error {
			var err error
			tx, err = userWalletClient.RecordTransaction(draftTx.GetDraftTransactionHex(), draftTx.GetDraftTransactionID(), metadata)
			return err //nolint:wrapcheck // error wrapped higher in call stack
		},
		retry.Attempts(retries),
		retry.Delay(1*time.Second),
		retry.OnRetry(func(n uint, err error) {
			log.Warn().
				Str("draftTxID", draftTx.GetDraftTransactionID()).
				Msgf("%d retry RecordTransaction after error: %v", n, err.Error())
		}),
	)
	return tx, err //nolint:wrapcheck // error wrapped higher in call stack
}
