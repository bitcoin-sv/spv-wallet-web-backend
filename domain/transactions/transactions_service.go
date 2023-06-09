package transactions

import (
	"bux-wallet/domain/users"
	"bux-wallet/logging"
	"fmt"

	"github.com/BuxOrg/go-buxclient/transports"
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
func (s *TransactionService) CreateTransaction(xpriv, recipient string, satoshis uint64) (string, error) {
	// Try to generate BUX client with decrypted xpriv.
	buxClient, err := s.buxClientFactory.CreateWithXpriv(xpriv)
	if err != nil {
		return "", err
	}

	fmt.Println(buxClient)

	// Create recipients.
	var recipients = []*transports.Recipients{
		{
			Satoshis: satoshis,
			To:       recipient,
		},
	}

	// Create transaction.
	transaction, err := buxClient.SendToRecipents(recipients)
	if err != nil {
		return "", err
	}

	return transaction, nil
}
