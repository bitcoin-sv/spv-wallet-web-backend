package websockets

import (
	"bux-wallet/logging"
	buxclient "bux-wallet/transports/bux/client"
	"encoding/json"
	"fmt"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/centrifugal/centrifuge"
	"time"
)

// Socket represents websocket server entrypoint used to publish messages via websocket communication.
type Socket struct {
	Client *centrifuge.Client
	Log    logging.Logger
}

// BaseEvent represents base of notification.
type BaseEvent struct {
	Status    string  `json:"status"`
	Error     *string `json:"error"`
	EventType string  `json:"eventType"`
}

// NewTransactionEvent represents notification about new transaction.
type NewTransactionEvent struct {
	BaseEvent
	Transaction transaction `json:"transaction"`
}

// transaction represents simplified transaction which is return in webhook.
type transaction struct {
	Id         string    `json:"id"`
	Receiver   string    `json:"receiver"`
	Sender     string    `json:"sender"`
	Status     string    `json:"status"`
	Direction  string    `json:"direction"`
	TotalValue uint64    `json:"totalValue"`
	CreatedAt  time.Time `json:"createdAt"`
}

// Notify send event notification.
func (s *Socket) Notify(event any) {
	bytes, err := json.Marshal(event)
	if err != nil {
		return
	}

	if err = s.Client.Send(bytes); err != nil {
		s.Log.Errorf("Error when sending event %v to websocket: %v", event, err.Error())
	}
	s.Log.Infof("Event %v sent to websocket", event)
}

// NotifyAboutTransaction will send notification about new transaction.
func (s *Socket) NotifyAboutTransaction(tx *buxmodels.Transaction) {
	txEvent := prepareNewTransactionEvent(tx)
	s.Notify(txEvent)
}

// SendError send error notification.
func (s *Socket) SendError(error error) {
	bytes, err := json.Marshal(error)
	if err != nil {
		return
	}

	if err = s.Client.Send(bytes); err != nil {
		s.Log.Errorf("Error when sending event %v to websocket: %v", error.Error(), err.Error())
	}
	s.Log.Infof("Event %v sent to websocket", error.Error())
}

func prepareNewTransactionEvent(tx *buxmodels.Transaction) NewTransactionEvent {
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
		Transaction: transaction{
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
