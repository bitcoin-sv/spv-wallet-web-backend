package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

// UpsertContact represents a request for creating or updating new contact.
type UpsertContact struct {
	FullName string         `json:"fullName"`
	Metadata map[string]any `json:"metadata"`
}

type SearchContact struct {
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	Metadata    models.Metadata        `json:"metadata,omitempty"`
	QueryParams *filter.QueryParams    `json:"params,omitempty"`
}

// ConfirmContact represents a request for confirming a contact.
type ConfirmContact struct {
	Passcode string          `json:"passcode"`
	Contact  *models.Contact `json:"contact,omitempty"`
}

// TotpResponse represents a response with generated passcode.
type TotpResponse struct {
	Passcode string `json:"passcode"`
}
