package transactions_test

import (
	"web-backend/notification"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"errors"
	"testing"

	"web-backend/tests/data"
	mock "web-backend/tests/mocks"
	"web-backend/tests/utils"

	"web-backend/domain/transactions"
	"web-backend/domain/users"
	"web-backend/transports/client"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	testLogger := zerolog.Nop()
	t.Run("Create transaction", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		paymail := "paymail@example.com"
		xpriv := gofakeit.HexUint256()
		recipient := "recipient.paymail@example.com"
		txValueInSatoshis := uint64(500)

		tr := client.DraftTransaction{}

		mockUserClient := mock.NewMockUserClient(ctrl)
		mockUserClient.EXPECT().
			CreateAndFinalizeTransaction(gomock.Any(), gomock.Any()).
			Return(&tr, nil)

		mockUserClient.EXPECT().
			RecordTransaction(gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		clientFctrMq := mock.NewMockClientFactory(ctrl)
		clientFctrMq.EXPECT().
			CreateWithXpriv(xpriv).
			Return(mockUserClient, nil)

		sut := transactions.NewTransactionService(mock.NewMockAdminClient(ctrl), clientFctrMq, &testLogger)

		// Act
		txs := make(chan notification.TransactionEvent, 1)
		err := sut.CreateTransaction(paymail, xpriv, recipient, txValueInSatoshis, txs)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetTransaction_ReturnsTransactionDetails(t *testing.T) {
	testLogger := zerolog.Nop()
	ts := data.CreateTestTransactions(10)

	cases := []struct {
		name          string
		transactionId string
	}{
		{
			name:          "Get transaction, return transaction details",
			transactionId: ts[0].GetTransactionId(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			paymail := "paymail@example.com"
			accessKey := gofakeit.HexUint256()

			mockUserClient := mock.NewMockUserClient(ctrl)
			mockUserClient.EXPECT().
				GetTransaction(tc.transactionId, paymail).
				Return(findById(ts, tc.transactionId))

			clientFctrMq := mock.NewMockClientFactory(ctrl)
			clientFctrMq.EXPECT().
				CreateWithAccessKey(accessKey).
				Return(mockUserClient, nil)

			sut := transactions.NewTransactionService(mock.NewMockAdminClient(ctrl), clientFctrMq, &testLogger)

			// Act
			result, err := sut.GetTransaction(accessKey, tc.transactionId, paymail)
			if err != nil {
				t.Fatal(err)
			}

			// Assert
			assert.Equal(t, tc.transactionId, result.GetTransactionId())
		})
	}
}

func TestGetTransaction_ReturnsError(t *testing.T) {
	testLogger := zerolog.Nop()
	ts := data.CreateTestTransactions(10)

	cases := []struct {
		name          string
		transactionId string
		expectdErr    error
	}{
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

			paymail := "paymail@example.com"
			accessKey := gofakeit.HexUint256()

			mockUserClient := mock.NewMockUserClient(ctrl)
			mockUserClient.EXPECT().
				GetTransaction(tc.transactionId, paymail).
				Return(findById(ts, tc.transactionId))

			clientFctrMq := mock.NewMockClientFactory(ctrl)
			clientFctrMq.EXPECT().
				CreateWithAccessKey(accessKey).
				Return(mockUserClient, nil)

			sut := transactions.NewTransactionService(mock.NewMockAdminClient(ctrl), clientFctrMq, &testLogger)

			// Act
			result, err := sut.GetTransaction(accessKey, tc.transactionId, paymail)

			// Assert
			require.EqualError(t, tc.expectdErr, err.Error())
			assert.Nil(t, result)
		})
	}
}

func findById(collection []client.FullTransaction, id string) (users.FullTransaction, error) {
	result := utils.Find(collection, func(t client.FullTransaction) bool { return t.Id == id })

	if result == nil {
		return nil, errors.New("Not found")
	}

	return result, nil
}
