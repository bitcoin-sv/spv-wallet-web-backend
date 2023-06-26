package transactions

import (
	"bux-wallet/domain/users"
	"bux-wallet/logging"
	"math"

	"github.com/BuxOrg/go-buxclient/transports"
	"github.com/mrz1836/go-datastore"
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
func (s *TransactionService) CreateTransaction(userPaymail, xpriv, recipient string, satoshis uint64) (users.Transaction, error) {
	// Try to generate BUX client with decrypted xpriv.
	buxClient, err := s.buxClientFactory.CreateWithXpriv(xpriv)
	if err != nil {
		return nil, err
	}

	// Create recipients.
	var recipients = []*transports.Recipients{
		{
			Satoshis: satoshis,
			To:       recipient,
		},
	}

	// Send transaction.
	transaction, err := buxClient.SendToRecipents(recipients, userPaymail)
	if err != nil {
		return nil, err
	}

	return transaction, nil
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
func (s *TransactionService) GetTransactions(accessKey, userPaymail string, queryParam datastore.QueryParams) (*PaginatedTransactions, error) {
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
