package data

import (
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/spvwallet"

	"github.com/brianvoe/gofakeit/v6"
)

// CreateTestTransactions returns 'count' randomly generated transactions.
func CreateTestTransactions(count int) []spvwallet.FullTransaction {
	result := make([]spvwallet.FullTransaction, count)
	gofakeit.Slice(&result)

	return result
}
