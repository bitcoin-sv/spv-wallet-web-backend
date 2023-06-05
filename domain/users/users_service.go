package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"bux-wallet/encryption"
	"bux-wallet/logging"
	buxclient "bux-wallet/transports/bux/client"

	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/bip39"
	"github.com/libsv/go-bk/chaincfg"
)

// UserService represents User service and provide access to repository.
type UserService struct {
	repo      UsersRepository
	BuxClient *buxclient.AdminBuxClient
	log       logging.Logger
}

// NewUserService creates UserService instance.
func NewUserService(repo UsersRepository, buxClient *buxclient.AdminBuxClient, lf logging.LoggerFactory) *UserService {
	return &UserService{
		repo:      repo,
		BuxClient: buxClient,
		log:       lf.NewLogger("user-service"),
	}
}

// InsertUser inserts user to database.
func (s *UserService) InsertUser(user *User) error {
	err := s.repo.InsertUser(context.Background(), user)
	return err
}

// CreateNewUser creates new user.
func (s *UserService) CreateNewUser(email, password string) (*CreatedUser, error) {
	// Validate password.
	err := validatePassword(password)
	if err != nil {
		return nil, err
	}

	// Validate user.
	err = s.validateUser(email)
	if err != nil {
		return nil, err
	}

	// Generate mnemonic and seed
	mnemonic, seed, err := generateMnemonic()
	if err != nil {
		return nil, err
	}

	xpriv, err := generateXpriv(seed)
	if err != nil {
		return nil, err
	}

	// Encrypt xpriv
	encryptedXpriv, err := encryptXpriv(password, xpriv.String())

	if err != nil {
		return nil, err
	}

	// Register xpub in BUX.
	xpub, err := s.BuxClient.RegisterXpub(xpriv)
	if err != nil {
		return nil, fmt.Errorf("error registering xpub in BUX: %s", err.Error())
	}

	// Get username from email which will be used as paymail alias.
	username, _ := splitEmail(email)

	// Register paymail in BUX.
	paymail, err := s.BuxClient.RegisterNewPaymail(username, xpub)
	if err != nil {
		return nil, fmt.Errorf("error registering paymail in BUX: %s", err.Error())
	}

	// Create and save new user.
	user := &User{
		Email:     email,
		Xpriv:     encryptedXpriv,
		Paymail:   paymail,
		CreatedAt: time.Now(),
	}

	err = s.InsertUser(user)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	// Decrypt xpriv.
	decryptedXpriv, err := decryptXpriv(password, user.Xpriv)
	if err != nil {
		return nil, err
	}

	// Try to generate BUX client with decrypted xpriv.
	buxClient, err := buxclient.CreateBuxClientFromRawXpriv(decryptedXpriv)
	if err != nil {
		return nil, err
	}

	// Create access key.
	accessKey, err := buxClient.CreateAccessKey()
	if err != nil {
		return nil, err
	}

	signInUser := &AuthenticatedUser{
		User: user,
		AccessKey: AccessKey{
			Id:  accessKey.Id,
			Key: accessKey.Key,
		},
	}

	return signInUser, nil
}

// SignOutUser signs out user by removing session and access key.
func (s *UserService) SignOutUser(accessKeyId string) error {
	// TODO: Revoke access key.
	//
	// err := s.BuxClient.RevokeAccessKey(accessKeyId)
	// if err != nil {
	// 	return err
	// }
	return nil
}

// GetUserById returns user by id.
func (s *UserService) GetUserById(userId int) (*User, error) {
	user, err := s.repo.GetUserById(context.Background(), userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) validateUser(email string) error {
	//Validate email
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email address")
	}

	// Check if user with email already exists.
	_, err = s.repo.GetUserByEmail(context.Background(), email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	return fmt.Errorf("user with email %s already exists", email)
}

// generateMnemonic generates mnemonic and seed.
func generateMnemonic() (string, []byte, error) {
	entropy, err := bip39.GenerateEntropy(160)
	if err != nil {
		return "", nil, err
	}

	return bip39.Mnemonic(entropy, "")
}

// generateXpriv generates xpriv from seed.
func generateXpriv(seed []byte) (*bip32.ExtendedKey, error) {
	xpriv, err := bip32.NewMaster(seed, &chaincfg.MainNet)
	if err != nil {
		return nil, err
	}
	return xpriv, nil
}

// encryptXpriv encrypts xpriv with password.
func encryptXpriv(password, xpriv string) (string, error) {
	// Create hash from password
	hashedPassword, err := encryption.Hash(password)
	if err != nil {
		return "", err
	}

	// Encrypt xpriv with hashed password
	encryptedXpriv, err := encryption.Encrypt(hashedPassword, xpriv)
	if err != nil {
		return "", err
	}

	return encryptedXpriv, nil
}

// decryptXpriv decrypts xpriv with password.
func decryptXpriv(password, encryptedXpriv string) (string, error) {
	// Create hash from password
	hashedPassword, err := encryption.Hash(password)
	if err != nil {
		return "", err
	}

	// Decrypt xpriv with hashed password
	xpriv := encryption.Decrypt(hashedPassword, encryptedXpriv)
	if err != nil {
		return "", err
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
