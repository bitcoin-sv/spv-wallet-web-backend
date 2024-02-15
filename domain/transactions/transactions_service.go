package transactions

import (
	"math"
	"time"

	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet-web-backend/notification"

	"github.com/rs/zerolog"

	models "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/go-buxclient/transports"
	"github.com/avast/retry-go/v4"
)

// TransactionService represents service whoch contains methods linked with transactions.
type TransactionService struct {
	adminClient   users.AdminClient
	clientFactory users.ClientFactory
	log           *zerolog.Logger
}

// NewTransactionService creates new transaction service.
func NewTransactionService(adminClient users.AdminClient, clientFactory users.ClientFactory, log *zerolog.Logger) *TransactionService {
	transactionServiceLogger := log.With().Str("service", "user-service").Logger()
	return &TransactionService{
		adminClient:   adminClient,
		clientFactory: clientFactory,
		log:           &transactionServiceLogger,
	}
}

// CreateTransaction creates transaction.
func (s *TransactionService) CreateTransaction(userPaymail, xpriv, recipient string, satoshis uint64, events chan notification.TransactionEvent) error {
	userClient, err := s.clientFactory.CreateWithXpriv(xpriv)
	if err != nil {
		return err
	}

	var recipients = []*transports.Recipients{
		{
			Satoshis: satoshis,
			To:       recipient,
		},
	}

	metadata := &models.Metadata{
		"receiver": recipient,
		"sender":   userPaymail,
	}

	draftTransaction, err := userClient.CreateAndFinalizeTransaction(recipients, metadata)
	if err != nil {
		return err
	}

	go func() {
		tx, err := tryRecordTransaction(userClient, draftTransaction, metadata, s.log)
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
	userClient, err := s.clientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return nil, err
	}

	transaction, err := userClient.GetTransaction(id, userPaymail)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransactions returns transactions by access key.
func (s *TransactionService) GetTransactions(accessKey, userPaymail string, queryParam transports.QueryParams) (*PaginatedTransactions, error) {
	// Try to generate user-client with decrypted xpriv.
	userClient, err := s.clientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return nil, err
	}

	count, err := userClient.GetTransactionsCount()
	if err != nil {
		return nil, err
	}

	transactions, err := userClient.GetTransactions(queryParam, userPaymail)
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

func tryRecordTransaction(userClient users.UserClient, draftTx users.DraftTransaction, metadata *models.Metadata, log *zerolog.Logger) (*models.Transaction, error) {
	retries := uint(3)
	tx, recordErr := tryRecord(userClient, draftTx, metadata, log, retries)

	// unreserve utxos if failed
	if recordErr != nil {
		log.Error().
			Str("draftTxId", draftTx.GetDraftTransactionId()).
			Msgf("record transaction failed: %s", recordErr.Error())

		unreserveErr := tryUnreserve(userClient, draftTx, log, retries)
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

func tryRecord(userClient users.UserClient, draftTx users.DraftTransaction, metadata *models.Metadata, log *zerolog.Logger, retries uint) (*models.Transaction, error) {
	log.Debug().
		Str("draftTxId", draftTx.GetDraftTransactionId()).
		Msg("record transaction")

	tx := &models.Transaction{}
	err := retry.Do(
		func() error {
			var err error
			tx, err = userClient.RecordTransaction(draftTx.GetDraftTransactionHex(), draftTx.GetDraftTransactionId(), metadata)
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

func tryUnreserve(userClient users.UserClient, draftTx users.DraftTransaction, log *zerolog.Logger, retries uint) error {
	log.Debug().
		Str("draftTxId", draftTx.GetDraftTransactionId()).
		Msg("unreserve UTXOs from draft")

	return retry.Do(
		func() error {
			return userClient.UnreserveUtxos(draftTx.GetDraftTransactionId())
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
