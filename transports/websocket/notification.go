package websocket

import (
	"bux-wallet/logging"
	"bux-wallet/notification"
	"encoding/json"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/centrifugal/centrifuge"
)

// Socket represents websocket server entrypoint used to publish messages via websocket communication.
type Socket struct {
	Client *centrifuge.Client
	Log    logging.Logger
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
	txEvent := notification.PrepareNewTransactionEvent(tx)
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
