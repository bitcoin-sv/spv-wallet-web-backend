package notification

import (
	buxclient "bux-wallet/transports/bux/client"
	"fmt"
	buxmodels "github.com/BuxOrg/bux-models"
	"time"
)

// BaseEvent represents base of notification.
type BaseEvent struct {
	Status    string  `json:"status"`
	Error     *string `json:"error"`
	EventType string  `json:"eventType"`
}

// NewTransactionEvent represents notification about new transaction.
type NewTransactionEvent struct {
	BaseEvent
	Transaction Transaction `json:"transaction"`
}

// Transaction represents simplified transaction which is return in webhook.
type Transaction struct {
	Id         string    `json:"id"`
	Receiver   string    `json:"receiver"`
	Sender     string    `json:"sender"`
	Status     string    `json:"status"`
	Direction  string    `json:"direction"`
	TotalValue uint64    `json:"totalValue"`
	CreatedAt  time.Time `json:"createdAt"`
}

func PrepareNewTransactionEvent(tx *buxmodels.Transaction) NewTransactionEvent {
	sender, receiver := buxclient.GetPaymailsFromMetadata(tx, "unknown")
	status := "unconfirmed"
	if tx.BlockHeight > 0 {
		status = "confirmed"
	}
	return NewTransactionEvent{
		BaseEvent: BaseEvent{
			Status:    "success",
			Error:     nil,
			EventType: "create_transaction",
		},
		Transaction: Transaction{
			Id:         tx.ID,
			Receiver:   receiver,
			Sender:     sender,
			Status:     status,
			Direction:  fmt.Sprint(tx.TransactionDirection),
			TotalValue: tx.TotalValue,
			CreatedAt:  tx.Model.CreatedAt,
		},
	}
}

// PrepareNewTransactionErrorEvent prepares error event in NewTransactionEvent struct.
func PrepareNewTransactionErrorEvent(err error) NewTransactionEvent {
	errString := err.Error()
	return NewTransactionEvent{
		BaseEvent: BaseEvent{
			Status:    "error",
			Error:     &errString,
			EventType: "create_transaction",
		},
		Transaction: Transaction{},
	}
}
