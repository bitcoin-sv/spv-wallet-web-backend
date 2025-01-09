package transactions_test

import (
	"errors"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/transactions"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet-web-backend/notification"
	"github.com/bitcoin-sv/spv-wallet-web-backend/spverrors"
	"github.com/bitcoin-sv/spv-wallet-web-backend/tests/data"
	mock "github.com/bitcoin-sv/spv-wallet-web-backend/tests/mocks"
	"github.com/bitcoin-sv/spv-wallet-web-backend/tests/utils"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/spvwallet"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		tr := spvwallet.DraftTransaction{}

		mockUserWalletClient := mock.NewMockUserWalletClient(ctrl)
		mockUserWalletClient.EXPECT().
			CreateAndFinalizeTransaction(gomock.Any(), gomock.Any()).
			Return(&tr, nil)

		mockUserWalletClient.EXPECT().
			RecordTransaction(gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()

		clientFctrMq := mock.NewMockWalletClientFactory(ctrl)
		clientFctrMq.EXPECT().
			CreateWithXpriv(xpriv).
			Return(mockUserWalletClient, nil)

		sut := transactions.NewTransactionService(mock.NewMockAdminWalletClient(ctrl), clientFctrMq, &testLogger)

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
		transactionID string
	}{
		{
			name:          "Get transaction, return transaction details",
			transactionID: ts[0].GetTransactionID(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			paymail := "paymail@example.com"
			accessKey := gofakeit.HexUint256()

			mockUserWalletClient := mock.NewMockUserWalletClient(ctrl)
			mockUserWalletClient.EXPECT().
				GetTransaction(tc.transactionID, paymail).
				Return(findByID(ts, tc.transactionID))

			clientFctrMq := mock.NewMockWalletClientFactory(ctrl)
			clientFctrMq.EXPECT().
				CreateWithAccessKey(accessKey).
				Return(mockUserWalletClient, nil)

			sut := transactions.NewTransactionService(mock.NewMockAdminWalletClient(ctrl), clientFctrMq, &testLogger)

			// Act
			result, err := sut.GetTransaction(accessKey, tc.transactionID, paymail)
			if err != nil {
				t.Fatal(err)
			}

			// Assert
			assert.Equal(t, tc.transactionID, result.GetTransactionID())
		})
	}
}

func TestGetTransaction_ReturnsError(t *testing.T) {
	testLogger := zerolog.Nop()
	ts := data.CreateTestTransactions(10)

	cases := []struct {
		name          string
		transactionID string
		expectdErr    error
	}{
		{
			name:          "Transaction doesn't exist",
			transactionID: "imnothere",
			expectdErr:    spverrors.ErrGetTransaction,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			paymail := "paymail@example.com"
			accessKey := gofakeit.HexUint256()

			mockUserWalletClient := mock.NewMockUserWalletClient(ctrl)
			mockUserWalletClient.EXPECT().
				GetTransaction(tc.transactionID, paymail).
				Return(findByID(ts, tc.transactionID))

			clientFctrMq := mock.NewMockWalletClientFactory(ctrl)
			clientFctrMq.EXPECT().
				CreateWithAccessKey(accessKey).
				Return(mockUserWalletClient, nil)

			sut := transactions.NewTransactionService(mock.NewMockAdminWalletClient(ctrl), clientFctrMq, &testLogger)

			// Act
			result, err := sut.GetTransaction(accessKey, tc.transactionID, paymail)

			// Assert
			require.EqualError(t, tc.expectdErr, err.Error())
			assert.Nil(t, result)
		})
	}
}

func findByID(collection []spvwallet.FullTransaction, id string) (users.FullTransaction, error) {
	result := utils.Find(collection, func(t spvwallet.FullTransaction) bool { return t.ID == id })

	if result == nil {
		return nil, errors.New("Not found")
	}

	return result, nil
}
