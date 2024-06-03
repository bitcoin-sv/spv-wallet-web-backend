package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

// CreateTransaction represents request for creating new transaction.
type CreateTransaction struct {
	Password  string `json:"password"`
	Recipient string `json:"recipient"`
	Satoshis  uint64 `json:"satoshis"`
}

type SearchTransaction struct {
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	Metadata    models.Metadata        `json:"metadata,omitempty"`
	QueryParams *filter.QueryParams    `json:"params,omitempty"`
}
