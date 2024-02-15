package transactions

import (
	"math"
	"time"

	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet-web-backend/notification"

	"github.com/rs/zerolog"

	walletmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/go-buxclient/transports"
	"github.com/avast/retry-go/v4"
)

// TransactionService represents service whoch contains methods linked with transactions.
type TransactionService struct {
	adminWalletClient   users.AdminWalletClient
	walletClientFactory users.WalletClientFactory
	log                 *zerolog.Logger
}

// NewTransactionService creates new transaction service.
func NewTransactionService(adminWalletClient users.AdminWalletClient, walletClientFactory users.WalletClientFactory, log *zerolog.Logger) *TransactionService {
	transactionServiceLogger := log.With().Str("service", "user-service").Logger()
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
		return err
	}

	var recipients = []*transports.Recipients{
		{
			Satoshis: satoshis,
			To:       recipient,
		},
	}

	metadata := &walletmodels.Metadata{
		"receiver": recipient,
		"sender":   userPaymail,
	}

	draftTransaction, err := userWalletClient.CreateAndFinalizeTransaction(recipients, metadata)
	if err != nil {
		return err
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
		return nil, err
	}

	transaction, err := userWalletClient.GetTransaction(id, userPaymail)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransactions returns transactions by access key.
func (s *TransactionService) GetTransactions(accessKey, userPaymail string, queryParam transports.QueryParams) (*PaginatedTransactions, error) {
	// Try to generate user-client with decrypted xpriv.
	userWalletClient, err := s.walletClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return nil, err
	}

	count, err := userWalletClient.GetTransactionsCount()
	if err != nil {
		return nil, err
	}

	transactions, err := userWalletClient.GetTransactions(queryParam, userPaymail)
	if err != nil {
		return nil, err
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

func tryRecordTransaction(userWalletClient users.UserWalletClient, draftTx users.DraftTransaction, metadata *walletmodels.Metadata, log *zerolog.Logger) (*walletmodels.Transaction, error) {
	retries := uint(3)
	tx, recordErr := tryRecord(userWalletClient, draftTx, metadata, log, retries)

	// unreserve utxos if failed
	if recordErr != nil {
		log.Error().
			Str("draftTxId", draftTx.GetDraftTransactionId()).
			Msgf("record transaction failed: %s", recordErr.Error())

		unreserveErr := tryUnreserve(userWalletClient, draftTx, log, retries)
		if unreserveErr != nil {
			log.Error().
				Str("draftTxId", draftTx.GetDraftTransactionId()).
				Msgf("unreserve transaction failed: %s", unreserveErr.Error())
		}
		return nil, recordErr
	}

	log.Debug().
		Str("draftTxId", draftTx.GetDraftTransactionId()).
		Msg("transaction successfully recorded")
	return tx, nil
}

func tryRecord(userWalletClient users.UserWalletClient, draftTx users.DraftTransaction, metadata *walletmodels.Metadata, log *zerolog.Logger, retries uint) (*walletmodels.Transaction, error) {
	log.Debug().
		Str("draftTxId", draftTx.GetDraftTransactionId()).
		Msg("record transaction")

	tx := &walletmodels.Transaction{}
	err := retry.Do(
		func() error {
			var err error
			tx, err = userWalletClient.RecordTransaction(draftTx.GetDraftTransactionHex(), draftTx.GetDraftTransactionId(), metadata)
			return err
		},
		retry.Attempts(retries),
		retry.Delay(1*time.Second),
		retry.OnRetry(func(n uint, err error) {
			log.Warn().
				Str("draftTxId", draftTx.GetDraftTransactionId()).
				Msgf("%d retry RecordTransaction after error: %v", n, err.Error())
		}),
	)
	return tx, err
}

func tryUnreserve(userWalletClient users.UserWalletClient, draftTx users.DraftTransaction, log *zerolog.Logger, retries uint) error {
	log.Debug().
		Str("draftTxId", draftTx.GetDraftTransactionId()).
		Msg("unreserve UTXOs from draft")

	return retry.Do(
		func() error {
			return userWalletClient.UnreserveUtxos(draftTx.GetDraftTransactionId())
		},
		retry.Attempts(retries),
		retry.Delay(1*time.Second),
		retry.OnRetry(func(n uint, err error) {
			log.Warn().
				Str("draftTxId", draftTx.GetDraftTransactionId()).
				Msgf("%d retry UnreserveUtxos after error: %v", n, err.Error())
		}),
	)
}
