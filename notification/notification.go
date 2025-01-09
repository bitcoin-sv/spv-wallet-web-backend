package notification

import (
	"fmt"
	"time"

	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/spvwallet"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// BaseEvent represents base of notification.
type BaseEvent struct {
	Status    string  `json:"status"`
	Error     *string `json:"error"`
	EventType string  `json:"eventType"`
}

// TransactionEvent represents notification about new transaction.
type TransactionEvent struct {
	BaseEvent
	Transaction *Transaction `json:"transaction"`
}

// Transaction represents simplified transaction which is return in webhook.
type Transaction struct {
	ID         string    `json:"id"`
	Receiver   string    `json:"receiver"`
	Sender     string    `json:"sender"`
	Status     string    `json:"status"`
	Direction  string    `json:"direction"`
	TotalValue uint64    `json:"totalValue"`
	CreatedAt  time.Time `json:"createdAt"`
}

// PrepareTransactionEvent prepares event in NewTransactionEvent struct.
func PrepareTransactionEvent(tx *models.Transaction) TransactionEvent {
	sender, receiver := spvwallet.GetPaymailsFromMetadata(&response.Transaction{
		Model:                response.Model(tx.Model),
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
	}, "unknown")
	status := "unconfirmed"
	if tx.BlockHeight > 0 {
		status = "confirmed"
	}
	return TransactionEvent{
		BaseEvent: BaseEvent{
			Status:    "success",
			Error:     nil,
			EventType: "create_transaction",
		},
		Transaction: &Transaction{
			ID:         tx.ID,
			Receiver:   receiver,
			Sender:     sender,
			Status:     status,
			Direction:  fmt.Sprint(tx.TransactionDirection),
			TotalValue: tx.TotalValue,
			CreatedAt:  tx.Model.CreatedAt,
		},
	}
}

// PrepareTransactionErrorEvent prepares error event in NewTransactionEvent struct.
func PrepareTransactionErrorEvent(err error) TransactionEvent {
	errString := err.Error()
	return TransactionEvent{
		BaseEvent: BaseEvent{
			Status:    "error",
			Error:     &errString,
			EventType: "create_transaction",
		},
		Transaction: nil,
	}
}
