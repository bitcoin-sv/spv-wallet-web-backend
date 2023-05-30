package users

import (
	"context"
	"fmt"
	"strings"
	"time"

	"bux-wallet/data/users"
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
func NewUserService(repo *users.UsersRepository, buxClient *buxclient.AdminBuxClient, lf logging.LoggerFactory) *UserService {
	return &UserService{
		repo:      repo,
		BuxClient: buxClient,
		log:       lf.NewLogger("user-service"),
	}
}

// InsertUser inserts user to database.
func (s *UserService) InsertUser(user *User) error {
	err := s.repo.InsertUser(context.Background(), user.toUserDto())
	return err
}

// CreateNewUser creates new user.
func (s *UserService) CreateNewUser(email, password string) (string, string, error) {
	// Check if user with email already exists.
	exists := s.checkIfUserExists(email)
	if exists {
		return "", "", fmt.Errorf("user with email %s already exists", email)
	}

	// Generate mnemonic and seed
	mnemonic, seed, err := generateMnemonic()
	if err != nil {
		return "", "", err
	}

	xpriv, err := generateXpriv(seed)
	if err != nil {
		return "", "", err
	}

	// Encrypt xpriv
	encryptedXpriv, err := encryptXpriv(password, xpriv.String())

	if err != nil {
		return "", "", err
	}

	// Register xpub in BUX.
	xpub, err := s.BuxClient.RegisterXpub(xpriv)
	if err != nil {
		return "", "", fmt.Errorf("error registering xpub in BUX: %s", err.Error())
	}

	// Get username from email which will be used as paymail alias.
	username, _ := splitEmail(email)

	// Register paymail in BUX.
	paymail, err := s.BuxClient.RegisterNewPaymail(username, xpub)
	if err != nil {
		return "", "", fmt.Errorf("error registering paymail in BUX: %s", err.Error())
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
		return "", "", err
	}

	return mnemonic, paymail, err
}

func (s *UserService) checkIfUserExists(email string) bool {
	_, err := s.repo.GetUserByEmail(context.Background(), email)
	return err == nil
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

// splitEmail splits email to username and domain.
func splitEmail(email string) (string, string) {
	components := strings.Split(email, "@")
	username, domain := components[0], components[1]

	return username, domain
}
