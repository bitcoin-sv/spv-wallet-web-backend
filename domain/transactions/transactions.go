package transactions

import "bux-wallet/domain/users"

type PaginatedTransactions struct {
	Count        int64               `json:"count"`
	Pages        int                 `json:"pages"`
	Transactions []users.Transaction `json:"transactions"`
}
