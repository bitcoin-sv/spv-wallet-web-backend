package contacts

import "github.com/bitcoin-sv/spv-wallet/models"

// UpsertContact represents a request for creating or updating new contact.
type UpsertContact struct {
	FullName string          `json:"fullName"`
	Metadata models.Metadata `json:"metadata"`
}

type SearchContact struct {
	Conditions map[string]interface{} `json:"conditions,omitempty"`
	Metadata   models.Metadata        `json:"metadata,omitempty"`
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
