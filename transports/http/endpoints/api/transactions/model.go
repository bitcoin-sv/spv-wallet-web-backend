package transactions

// CreateTransaction represents request for creating new transaction.
type CreateTransaction struct {
	Password  string `json:"password"`
	Recipient string `json:"recipent"`
	Satoshis  uint64 `json:"satoshis"`
}
