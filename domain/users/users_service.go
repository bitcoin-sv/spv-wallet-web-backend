package users

import (
	"context"
	"database/sql"
	"fmt"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/rates"
	"github.com/bitcoin-sv/spv-wallet-web-backend/encryption"
	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/bip39"
	"github.com/libsv/go-bk/chaincfg"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// CredentialsError Generic error type / wrapper for Credentials errors.
type CredentialsError struct {
	message string
}

func (e *CredentialsError) Error() string {
	return e.message
}

// UserError Generic error type / wrapper for User errors.
type UserError struct {
	message string
}

func (e *UserError) Error() string {
	return e.message
}

// PaymailError Generic error type / wrapper for Paymail errors.
type PaymailError struct {
	message string
}

func (e *PaymailError) Error() string {
	return e.message
}

// XPubError Generic error type / wrapper for XPub errors.
type XPubError struct {
	message string
}

func (e *XPubError) Error() string {
	return e.message
}

// ErrInvalidCredentials is throwing when invalid credentials were used.
var ErrInvalidCredentials = &CredentialsError{"invalid credentials"}

// ErrUserAlreadyExists is throwing when we try to register a user with already used email.
var ErrUserAlreadyExists = &UserError{"user already exists"}

// UserService represents User service and provide access to repository.
type UserService struct {
	repo                Repository
	ratesService        *rates.Service
	adminWalletClient   AdminWalletClient
	walletClientFactory WalletClientFactory
	log                 *zerolog.Logger
}

// NewUserService creates UserService instance.
func NewUserService(repo Repository, adminWalletClient AdminWalletClient, walletClientFactory WalletClientFactory, rService *rates.Service, l *zerolog.Logger) *UserService {
	userServiceLogger := l.With().Str("service", "user-service").Logger()
	s := &UserService{
		repo:                repo,
		adminWalletClient:   adminWalletClient,
		walletClientFactory: walletClientFactory,
		ratesService:        rService,
		log:                 &userServiceLogger,
	}

	return s
}

// InsertUser inserts user to database.
func (s *UserService) InsertUser(user *User) error {
	if err := s.repo.InsertUser(context.Background(), user); err != nil {
		e := &UserError{err.Error()}
		s.log.Error().Msgf("Error while inserting user: %v", e.Error())
		return e
	}
	return nil
}

// CreateNewUser creates new user.
func (s *UserService) CreateNewUser(email, password string) (*CreatedUser, error) {
	if err := validatePassword(password); err != nil {
		e := &UserError{err.Error()}
		s.log.Error().
			Str("newUserEmail", email).
			Msgf("Error while validating password: %v", e.Error())
		return nil, e
	}

	if err := s.validateUser(email); err != nil {
		e := &UserError{err.Error()}
		s.log.Error().
			Str("newUserEmail", email).
			Msgf("Error while validating user: %v", e.Error())
		return nil, e
	}

	mnemonic, seed, err := generateMnemonic()
	if err != nil {
		e := &UserError{err.Error()}
		s.log.Error().
			Str("newUserEmail", email).
			Msgf("Error while generating mnemonic: %v", e.Error())
		return nil, e
	}

	xpriv, err := generateXpriv(seed)
	if err != nil {
		e := &UserError{err.Error()}
		s.log.Error().
			Str("newUserEmail", email).
			Msgf("Error while generating xPriv: %v", e.Error())
		return nil, e
	}

	encryptedXpriv, err := encryptXpriv(password, xpriv.String())
	if err != nil {
		e := &UserError{err.Error()}
		s.log.Error().
			Str("newUserEmail", email).
			Msgf("Error while encrypting xPriv: %v", e.Error())
		return nil, e
	}

	xpub, err := s.adminWalletClient.RegisterXpub(xpriv)
	if err != nil {
		e := &XPubError{fmt.Sprintf("error registering xpub in spv-wallet: %s", err)}
		s.log.Error().
			Str("newUserEmail", email).
			Msgf("Error while registering xPub: %v", e.Error())
		return nil, e
	}

	username, _ := splitEmail(email)

	paymail, err := s.adminWalletClient.RegisterPaymail(username, xpub)
	if err != nil {
		e := &PaymailError{fmt.Sprintf("error registering paymail in spv-wallet: %s", err)}
		s.log.Error().
			Str("newUserEmail", email).
			Msgf("Error while registering paymail: %v", e.Error())
		return nil, e
	}

	user := &User{
		Email:     email,
		Xpriv:     encryptedXpriv,
		Paymail:   paymail,
		CreatedAt: time.Now(),
	}

	if err = s.InsertUser(user); err != nil {
		e := &UserError{err.Error()}
		s.log.Error().
			Str("newUserEmail", email).
			Msgf("Error while inserting user: %v", e.Error())
		return nil, e
	}

	newUSerData := &CreatedUser{
		User:     user,
		Mnemonic: mnemonic,
	}

	return newUSerData, err
}

// SignInUser signs in user.
func (s *UserService) SignInUser(email, password string) (*AuthenticatedUser, error) {
	// Check if user exists.
	user, err := s.repo.GetUserByEmail(context.Background(), email)
	if err != nil {
		s.log.Error().
			Str("userEmail", email).
			Msgf("User wasn't found by email: %v", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, errors.Wrap(err, "cannot get user")
	}

	// Decrypt xpriv.
	decryptedXpriv, err := decryptXpriv(password, user.Xpriv)
	if err != nil {
		s.log.Error().
			Str("userEmail", email).
			Msgf("Error while decrypting xPriv: %v", err.Error())
		return nil, errors.Wrap(err, "cannot get xpriv")
	}

	userWalletClient := s.walletClientFactory.CreateWithXpriv(decryptedXpriv)

	accessKey, err := userWalletClient.CreateAccessKey()
	if err != nil {
		s.log.Error().
			Str("userEmail", email).
			Msgf("Error while creating access key: %v", err.Error())
		return nil, errors.Wrap(err, "cannot create access key")
	}

	xpub, err := userWalletClient.GetXPub()
	if err != nil {
		s.log.Error().
			Str("userEmail", email).
			Msgf("Error while getting xPub: %v", err.Error())
		return nil, errors.Wrap(err, "cannot get spv wallet user info")
	}

	exchangeRate, err := s.ratesService.GetExchangeRate()
	if err != nil {
		s.log.Error().
			Str("userEmail", email).
			Msgf("Exchange rate not found: %v", err.Error())
		return nil, errors.Wrap(err, "cannot get exchange rate")
	}

	balance, err := calculateBalance(xpub.GetCurrentBalance(), exchangeRate)
	if err != nil {
		s.log.Error().
			Str("userEmail", email).
			Msgf("Error while calculating balance: %v", err.Error())
		return nil, errors.Wrap(err, "cannot calculate balance")
	}

	signInUser := &AuthenticatedUser{
		User: user,
		AccessKey: AccessKey{
			ID:  accessKey.GetAccessKeyID(),
			Key: accessKey.GetAccessKey(),
		},
		Balance: *balance,
		Xpriv:   decryptedXpriv,
	}

	return signInUser, nil
}

// SignOutUser signs out user by revoking access key. (Not possible at the moment, method is just a mock.)
func (s *UserService) SignOutUser(_, _ string) error {

	// / Right now we cannot revoke access key without authentication with XPriv, which is impossible here.

	// userWalletClient := s.walletClientFactory.CreateWithAccessKey(accessKey)

	// _, err = userWalletClient.RevokeAccessKey(accessKeyId)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// GetUserByID returns user by id.
func (s *UserService) GetUserByID(userID int) (*User, error) {
	user, err := s.repo.GetUserByID(context.Background(), userID)
	if err != nil {
		s.log.Error().
			Str("userID", strconv.Itoa(userID)).
			Msgf("Error while getting user by id: %v", err.Error())
		return nil, fmt.Errorf("cannot get user, cause: %w", err)
	}

	return user, nil
}

// GetUserBalance returns user balance using access key.
func (s *UserService) GetUserBalance(accessKey string) (*Balance, error) {
	userWalletClient := s.walletClientFactory.CreateWithAccessKey(accessKey)

	// Get xpub.
	xpub, err := userWalletClient.GetXPub()
	if err != nil {
		s.log.Error().
			Msgf("Error while getting xPub: %v", err.Error())
		return nil, fmt.Errorf("cannot get user balance, cause: %w", err)
	}

	exchangeRate, err := s.ratesService.GetExchangeRate()
	if err != nil {
		s.log.Error().
			Msgf("Exchange rate not found: %v", err.Error())
		return nil, fmt.Errorf("exchange rate not found, cause: %w", err)
	}

	// Calculate balance.
	balance, err := calculateBalance(xpub.GetCurrentBalance(), exchangeRate)
	if err != nil {
		s.log.Error().
			Str("xpubID", xpub.GetID()).
			Msgf("Error while calculating balance: %v", err.Error())
		return nil, err
	}

	return balance, nil
}

// GetUserXpriv gets user by id and decrypt xpriv.
func (s *UserService) GetUserXpriv(userID int, password string) (string, error) {
	user, err := s.repo.GetUserByID(context.Background(), userID)
	if err != nil {
		s.log.Error().
			Str("userID", strconv.Itoa(userID)).
			Msgf("Error while getting user by id: %v", err.Error())

		return "", fmt.Errorf("cannot get user, cause: %w", err)
	}

	// Decrypt xpriv.
	decryptedXpriv, err := decryptXpriv(password, user.Xpriv)
	if err != nil {
		s.log.Error().
			Str("userID", strconv.Itoa(userID)).
			Msgf("Error while decrypting xPriv: %v", err.Error())
		return "", fmt.Errorf("cannot decrypt xpriv, cause: %w", err)
	}

	return decryptedXpriv, nil
}

func (s *UserService) validateUser(email string) error {
	// Validate email
	if _, err := mail.ParseAddress(email); err != nil {
		e := &UserError{"invalid email address"}
		s.log.Error().
			Str("userEmail", email).
			Msgf("Error while validating email: %v", e.Error())
		return e
	}

	// Check if user with email already exists.
	if _, err := s.repo.GetUserByEmail(context.Background(), email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		s.log.Error().
			Str("userEmail", email).
			Msgf("Error while validating email: %v", err.Error())
		return &UserError{err.Error()}
	}

	return ErrUserAlreadyExists
}

// generateMnemonic generates mnemonic and seed.
func generateMnemonic() (string, []byte, error) {
	entropy, err := bip39.GenerateEntropy(160)
	if err != nil {
		return "", nil, err //nolint:wrapcheck // error wrapped higher in call stack
	}

	return bip39.Mnemonic(entropy, "") //nolint:wrapcheck // error wrapped higher in call stack
}

// generateXpriv generates xpriv from seed.
func generateXpriv(seed []byte) (*bip32.ExtendedKey, error) {
	xpriv, err := bip32.NewMaster(seed, &chaincfg.MainNet)
	if err != nil {
		return nil, err //nolint:wrapcheck // error wrapped higher in call stack
	}
	return xpriv, nil
}

// encryptXpriv encrypts xpriv with password.
func encryptXpriv(password, xpriv string) (string, error) {
	// Create hash from password
	hashedPassword, err := encryption.Hash(password)
	if err != nil {
		return "", err //nolint:wrapcheck // error wrapped higher in call stack
	}

	// Encrypt xpriv with hashed password
	encryptedXpriv, err := encryption.Encrypt(hashedPassword, xpriv)
	if err != nil {
		return "", err //nolint:wrapcheck // error wrapped higher in call stack
	}

	return encryptedXpriv, nil
}

// decryptXpriv decrypts xpriv with password.
func decryptXpriv(password, encryptedXpriv string) (string, error) {
	// Create hash from password
	hashedPassword, err := encryption.Hash(password)
	if err != nil {
		return "", fmt.Errorf("internal error: %w", err)
	}

	// Decrypt xpriv with hashed password
	xpriv := encryption.Decrypt(hashedPassword, encryptedXpriv)
	if xpriv == "" {
		return "", ErrInvalidCredentials
	}

	return xpriv, nil
}

// splitEmail splits email to username and domain.
func splitEmail(email string) (string, string) {
	components := strings.Split(email, "@")
	username, domain := components[0], components[1]

	return username, domain
}

// validatePassword trim and validates password.
func validatePassword(password string) error {
	trimedPassword := strings.TrimSpace(password)
	if trimedPassword == "" {
		return fmt.Errorf("correct password is required")
	}

	return nil
}

func calculateBalance(satoshis uint64, exchangeRate *float64) (*Balance, error) {
	balanceBSV := float64(satoshis) / 100000000
	balanceUSD := balanceBSV * *exchangeRate

	balance := &Balance{
		Bsv:      balanceBSV,
		Usd:      balanceUSD,
		Satoshis: satoshis,
	}

	return balance, nil
}
