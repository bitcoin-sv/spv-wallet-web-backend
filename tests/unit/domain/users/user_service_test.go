package users_test

import (
	"database/sql"
	"errors"
	"testing"

	"bux-wallet/domain/users"
	"bux-wallet/logging"
	mock "bux-wallet/tests/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewUser(t *testing.T) {
	cases := []struct {
		name         string
		userEmail    string
		userPswd     string
		expectedUser *users.CreatedUser
		expectedErr  error
	}{
		{
			name:      "Insert valid user",
			userEmail: "homer.simpson@4chain.com",
			userPswd:  "strongP4$$word",
			expectedUser: &users.CreatedUser{
				User: &users.User{
					Email:   "homer.simpson@4chain.com",
					Paymail: "homer.simpson@homer.simpson.space",
				},
			},
			expectedErr: nil,
		},
		{
			name:         "User already exists",
			userEmail:    "marge.simpson@4chain.com",
			userPswd:     "strongP4$$word",
			expectedUser: nil,
			expectedErr:  errors.New("user with email marge.simpson@4chain.com already exists"),
		},
		{
			name:         "Invalid password",
			userEmail:    "ned.flanders@4chain.com",
			userPswd:     "",
			expectedUser: nil,
			expectedErr:  errors.New("correct password is required"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMq := mock.NewMockUsersRepository(ctrl)
			buxClientMq := mock.NewMockAdmBuxClient(ctrl)

			if tc.expectedUser != nil {
				repoMq.EXPECT().
					GetUserByEmail(gomock.Any(), tc.userEmail).
					Return(nil, sql.ErrNoRows)

				repoMq.EXPECT().InsertUser(gomock.Any(), gomock.Any())

				buxClientMq.EXPECT().
					RegisterXpub(gomock.Any()).
					Return(gomock.Any().String(), nil)
				buxClientMq.EXPECT().
					RegisterPaymail(gomock.Any(), gomock.Any()).
					Return(tc.expectedUser.User.Paymail, nil)

			} else {
				repoMq.EXPECT().
					GetUserByEmail(gomock.Any(), tc.userEmail).
					Return(nil, nil).
					AnyTimes()
			}

			sut := users.NewUserService(repoMq, buxClientMq, nil, logging.DefaultLoggerFactory())

			// Act
			result, err := sut.CreateNewUser(tc.userEmail, tc.userPswd)

			// Assert
			if err == nil {
				assertNewUser(t, tc.expectedUser, result)
			} else {
				assert.EqualError(t, err, tc.expectedErr.Error())
			}
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
