package websocket

import (
	"encoding/json"

	"github.com/centrifugal/centrifuge"
	"github.com/rs/zerolog"
)

// Socket represents websocket server entrypoint used to publish messages via websocket communication.
type Socket struct {
	Client *centrifuge.Client
	Log    *zerolog.Logger
}

// Notify send event notification.
func (s *Socket) Notify(event any) {
	bytes, err := json.Marshal(event)
	if err != nil {
		s.Log.Error().Msgf("Error when marshalling event %v: %v", event, err.Error())
		return
	}

	if s.Client == nil {
		s.Log.Debug().Msgf("Skipping notification, no client connection to handle the event %s", bytes)
		return
	}

	if err = s.Client.Send(bytes); err != nil {
		s.Log.Error().Msgf("Error when sending event %v to websocket: %v", event, err.Error())
	}
	s.Log.Info().Msgf("Event %v sent to websocket", event)
}
