package users_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/rs/zerolog"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	mock "github.com/bitcoin-sv/spv-wallet-web-backend/tests/mocks"
)

func TestCreateNewUser_ReturnsUser(t *testing.T) {
	testLogger := zerolog.Nop()
	cases := []struct {
		name         string
		userEmail    string
		userPswd     string
		expectedUser *users.CreatedUser
	}{
		{
			name:      "Insert valid user",
			userEmail: "homer.simpson@example.com",
			userPswd:  "strongP4$$word",
			expectedUser: &users.CreatedUser{
				User: &users.User{
					Email:   "homer.simpson@example.com",
					Paymail: "homer.simpson@homer.simpson.space",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMq := mock.NewMockUsersRepository(ctrl)
			mockAdminClient := mock.NewMockAdminClient(ctrl)

			repoMq.EXPECT().
				GetUserByEmail(gomock.Any(), tc.userEmail).
				Return(nil, sql.ErrNoRows)

			repoMq.EXPECT().InsertUser(gomock.Any(), gomock.Any())

			mockAdminClient.EXPECT().
				RegisterXpub(gomock.Any()).
				Return(gomock.Any().String(), nil)
			mockAdminClient.EXPECT().
				RegisterPaymail(gomock.Any(), gomock.Any()).
				Return(tc.expectedUser.User.Paymail, nil)

			sut := users.NewUserService(repoMq, mockAdminClient, nil, &testLogger)

			// Act
			result, err := sut.CreateNewUser(tc.userEmail, tc.userPswd)
			if err != nil {
				t.Fatal(err)
			}

			// Assert
			assertNewUser(t, tc.expectedUser, result)
		})
	}
}

func TestCreateNewUser_InvalidData_ReturnsError(t *testing.T) {
	testLogger := zerolog.Nop()
	cases := []struct {
		name        string
		userEmail   string
		userPswd    string
		expectedErr error
	}{
		{
			name:        "User already exists",
			userEmail:   "marge.simpson@example.com",
			userPswd:    "strongP4$$word",
			expectedErr: users.ErrUserAlreadyExists,
		},
		{
			name:        "Invalid email",
			userEmail:   "bart.simpson_example.com",
			userPswd:    "strongP4$$word",
			expectedErr: errors.New("invalid email address"),
		},
		{
			name:        "Invalid password",
			userEmail:   "ned.flanders@example.com",
			userPswd:    "",
			expectedErr: errors.New("correct password is required"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMq := mock.NewMockUsersRepository(ctrl)
			mockAdminClient := mock.NewMockAdminClient(ctrl)

			repoMq.EXPECT().
				GetUserByEmail(gomock.Any(), tc.userEmail).
				Return(nil, nil).
				AnyTimes()

			sut := users.NewUserService(repoMq, mockAdminClient, nil, &testLogger)

			// Act
			result, err := sut.CreateNewUser(tc.userEmail, tc.userPswd)

			// Assert
			require.EqualError(t, err, tc.expectedErr.Error())
			assert.Nil(t, result)
		})
	}
}

func assertNewUser(t *testing.T, expectedUser, newUser *users.CreatedUser) {
	assert.Equal(t, expectedUser.User.Email, newUser.User.Email)
	assert.Equal(t, expectedUser.User.Paymail, newUser.User.Paymail)
	assert.NotEmpty(t, newUser.User.Xpriv)
	assert.NotEmpty(t, newUser.User.CreatedAt)
	assert.NotEmpty(t, newUser.Mnemonic)
}
