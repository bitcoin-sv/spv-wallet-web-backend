package transactions

import (
	"bux-wallet/domain/users"
	"bux-wallet/logging"
	"bux-wallet/notification"
	"math"
	"time"

	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/go-buxclient/transports"
	"github.com/avast/retry-go/v4"
)

// TransactionService represents service whoch contains methods linked with transactions.
type TransactionService struct {
	buxClient        users.AdmBuxClient
	buxClientFactory users.BuxClientFactory
	log              logging.Logger
}

// NewTransactionService creates new transaction service.
func NewTransactionService(buxClient users.AdmBuxClient, bf users.BuxClientFactory, lf logging.LoggerFactory) *TransactionService {
	return &TransactionService{
		buxClient:        buxClient,
		buxClientFactory: bf,
		log:              lf.NewLogger("transaction-service"),
	}
}

// CreateTransaction creates transaction.
func (s *TransactionService) CreateTransaction(userPaymail, xpriv, recipient string, satoshis uint64, events chan notification.TransactionEvent) error {
	buxClient, err := s.buxClientFactory.CreateWithXpriv(xpriv)
	if err != nil {
		return err
	}

	var recipients = []*transports.Recipients{
		{
			Satoshis: satoshis,
			To:       recipient,
		},
	}

	metadata := &buxmodels.Metadata{
		"receiver": recipient,
		"sender":   userPaymail,
	}

	draftTransaction, err := buxClient.CreateAndFinalizeTransaction(recipients, metadata)
	if err != nil {
		return err
	}

	go func() {
		tx, err := tryRecordTransaction(buxClient, draftTransaction, metadata, s.log)
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
	// Try to generate BUX client with decrypted xpriv.
	buxClient, err := s.buxClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return nil, err
	}

	transaction, err := buxClient.GetTransaction(id, userPaymail)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransactions returns transactions by access key.
func (s *TransactionService) GetTransactions(accessKey, userPaymail string, queryParam transports.QueryParams) (*PaginatedTransactions, error) {
	// Try to generate BUX client with decrypted xpriv.
	buxClient, err := s.buxClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return nil, err
	}

	count, err := buxClient.GetTransactionsCount()
	if err != nil {
		return nil, err
	}

	transactions, err := buxClient.GetTransactions(queryParam, userPaymail)
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

func tryRecordTransaction(buxClient users.UserBuxClient, draftTx users.DraftTransaction, metadata *buxmodels.Metadata, log logging.Logger) (*buxmodels.Transaction, error) {
	retries := uint(3)
	tx, recordErr := tryRecord(buxClient, draftTx, metadata, log, retries)

	// unreserve utxos if failed
	if recordErr != nil {
		log.Errorf("record transaction failed: %s", recordErr.Error())

		unreserveErr := tryUnreserve(buxClient, draftTx, log, retries)
		if unreserveErr != nil {
			log.Errorf("unreserve transaction failed: %s", unreserveErr.Error())
		}
		return nil, recordErr
	}

	log.Debugf("transaction %s successfully recorded", draftTx.GetDraftTransactionId())
	return tx, nil
}

func tryRecord(buxClient users.UserBuxClient, draftTx users.DraftTransaction, metadata *buxmodels.Metadata, log logging.Logger, retries uint) (*buxmodels.Transaction, error) {
	log.Debugf("recording transaction %s", draftTx.GetDraftTransactionId())

	tx := &buxmodels.Transaction{}
	err := retry.Do(
		func() error {
			var err error
			tx, err = buxClient.RecordTransaction(draftTx.GetDraftTransactionHex(), draftTx.GetDraftTransactionId(), metadata)
			return err
		},
		retry.Attempts(retries),
		retry.Delay(1*time.Second),
		retry.OnRetry(func(n uint, err error) {
			log.Warnf("%d retry RecordTransaction after error: %s", n, err.Error())
		}),
	)
	return tx, err
}

func tryUnreserve(buxClient users.UserBuxClient, draftTx users.DraftTransaction, log logging.Logger, retries uint) error {
	log.Debugf("unreserve UTXOs from draft %s", draftTx.GetDraftTransactionId())

	return retry.Do(
		func() error {
			return buxClient.UnreserveUtxos(draftTx.GetDraftTransactionId())
		},
		retry.Attempts(retries),
		retry.Delay(1*time.Second),
		retry.OnRetry(func(n uint, err error) {
			log.Warnf("%d retry UnreserveUtxos after error: %s", n, err.Error())
		}),
	)
}
