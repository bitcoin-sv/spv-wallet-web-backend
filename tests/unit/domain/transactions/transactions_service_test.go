package transactions_test

import (
	"errors"
	"testing"

	"bux-wallet/tests/data"
	mock "bux-wallet/tests/mocks"
	"bux-wallet/tests/utils"

	"bux-wallet/domain/transactions"
	"bux-wallet/domain/users"
	"bux-wallet/logging"
	buxclient "bux-wallet/transports/bux/client"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	t.Run("Create transaction", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		paymail := "paymail@4chain.com"
		tr := buxclient.Transaction{}

		buxClientMq := mock.NewMockUserBuxClient(ctrl)
		buxClientMq.EXPECT().
			SendToRecipents(gomock.Any(), paymail).
			Return(&tr, nil)

		clientFctrMq := mock.NewMockBuxClientFactory(ctrl)
		clientFctrMq.EXPECT().
			CreateWithXpriv(gomock.Any()).
			Return(buxClientMq, nil)

		sut := transactions.NewTransactionService(mock.NewMockAdmBuxClient(ctrl), clientFctrMq, logging.DefaultLoggerFactory())

		// Act
		result, _ := sut.CreateTransaction(paymail, gomock.Any().String(), gomock.Any().String(), gofakeit.Uint64())

		// Assert
		assert.NotNil(t, result)
	})
}

func TestGetTransaction(t *testing.T) {
	ts := data.CreateTestTransactions(10)

	cases := []struct {
		name          string
		transactionId string
		expectdErr    error
	}{
		{
			name:          "Get transaction, return transaction details",
			transactionId: ts[0].GetTransactionId(),
		},
		{
			name:          "Transaction doesn't exist",
			transactionId: "imnothere",
			expectdErr:    errors.New("Not found"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			buxClientMq := mock.NewMockUserBuxClient(ctrl)
			buxClientMq.EXPECT().
				GetTransaction(tc.transactionId, gomock.Any().String()).
				Return(findById(ts, tc.transactionId))

			clientFctrMq := mock.NewMockBuxClientFactory(ctrl)
			clientFctrMq.EXPECT().
				CreateWithAccessKey(gomock.Any()).
				Return(buxClientMq, nil)

			sut := transactions.NewTransactionService(mock.NewMockAdmBuxClient(ctrl), clientFctrMq, logging.DefaultLoggerFactory())

			// Act
			result, err := sut.GetTransaction("fake-access-key", tc.transactionId, gomock.Any().String())

			// Assert
			if err != nil {
				assert.EqualError(t, tc.expectdErr, err.Error())
			} else {
				assert.Equal(t, tc.transactionId, result.GetTransactionId())
			}

		})
	}
}

func findById(collection []buxclient.FullTransaction, id string) (users.FullTransaction, error) {
	result := utils.Find(collection, func(t buxclient.FullTransaction) bool { return t.Id == id })

	if result == nil {
		return nil, errors.New("Not found")
	}

	return result, nil
}
