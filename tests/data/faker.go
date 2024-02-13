package data

import (
	"web-backend/transports/client"

	"github.com/brianvoe/gofakeit/v6"
)

// CreateTestTransactions returns 'count' randomly generated transactions.
func CreateTestTransactions(count int) []client.FullTransaction {
	result := make([]client.FullTransaction, count)
	gofakeit.Slice(&result)

	return result
}
