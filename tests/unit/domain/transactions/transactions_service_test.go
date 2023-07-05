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
		xpriv := gofakeit.HexUint256()
		recipient := "recipient.paymail@4chain.com"
		txValueInSatoshis := uint64(500)

		tr := buxclient.DraftTransaction{}

		buxClientMq := mock.NewMockUserBuxClient(ctrl)
		buxClientMq.EXPECT().
			CreateAndFinalizeTransaction(gomock.Any(), gomock.Any()).
			Return(&tr, nil)

		buxClientMq.EXPECT().
			RecordTransaction(gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		clientFctrMq := mock.NewMockBuxClientFactory(ctrl)
		clientFctrMq.EXPECT().
			CreateWithXpriv(xpriv).
			Return(buxClientMq, nil)

		sut := transactions.NewTransactionService(mock.NewMockAdmBuxClient(ctrl), clientFctrMq, logging.DefaultLoggerFactory())

		// Act
		err := sut.CreateTransaction(paymail, xpriv, recipient, txValueInSatoshis)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetTransaction_ReturnsTransactionDetails(t *testing.T) {
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

			paymail := "paymail@4chain.com"
			accessKey := gofakeit.HexUint256()

			buxClientMq := mock.NewMockUserBuxClient(ctrl)
			buxClientMq.EXPECT().
				GetTransaction(tc.transactionId, paymail).
				Return(findById(ts, tc.transactionId))

			clientFctrMq := mock.NewMockBuxClientFactory(ctrl)
			clientFctrMq.EXPECT().
				CreateWithAccessKey(accessKey).
				Return(buxClientMq, nil)

			sut := transactions.NewTransactionService(mock.NewMockAdmBuxClient(ctrl), clientFctrMq, logging.DefaultLoggerFactory())

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

			paymail := "paymail@4chain.com"
			accessKey := gofakeit.HexUint256()

			buxClientMq := mock.NewMockUserBuxClient(ctrl)
			buxClientMq.EXPECT().
				GetTransaction(tc.transactionId, paymail).
				Return(findById(ts, tc.transactionId))

			clientFctrMq := mock.NewMockBuxClientFactory(ctrl)
			clientFctrMq.EXPECT().
				CreateWithAccessKey(accessKey).
				Return(buxClientMq, nil)

			sut := transactions.NewTransactionService(mock.NewMockAdmBuxClient(ctrl), clientFctrMq, logging.DefaultLoggerFactory())

			// Act
			result, err := sut.GetTransaction(accessKey, tc.transactionId, paymail)

			// Assert
			assert.EqualError(t, tc.expectdErr, err.Error())
			assert.Nil(t, result)
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
