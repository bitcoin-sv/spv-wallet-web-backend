package transactions

// CreateTransaction represents request for creating new transaction.
type CreateTransaction struct {
	Password  string `json:"password"`
	Recipient string `json:"recipient"`
	Satoshis  uint64 `json:"satoshis"`
	Data      string `json:"data"`
}
