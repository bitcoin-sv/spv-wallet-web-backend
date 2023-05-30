package users

import (
	"bux-wallet/data/users"
	"bux-wallet/encryption"
	buxclient "bux-wallet/transports/bux/client"
	"context"
	"time"

	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/bip39"
	"github.com/libsv/go-bk/chaincfg"
)

// UserService represents User service and provide access to repository.
type UserService struct {
	repo      UsersRepository
	BuxClient *buxclient.AdminBuxClient
}

// NewUserService creates UserService instance.
func NewUserService(repo *users.UsersRepository, buxClient *buxclient.AdminBuxClient) *UserService {
	return &UserService{
		repo:      repo,
		BuxClient: buxClient,
	}
}

// InsertUser inserts user to database.
func (s *UserService) InsertUser(user *User) error {
	err := s.repo.InsertUser(context.Background(), user.toUserDto())
	return err
}

// CreateNewUser creates new user.
func (s *UserService) CreateNewUser(email, password string) (string, error) {
	// Generate mnemonic and seed
	mnemonic, seed, err := generateMnemonic()
	if err != nil {
		return "", err
	}

	xpriv, err := generateXpriv(seed)
	if err != nil {
		return "", err
	}

	// Create hash from password
	hashedPassword, err := encryption.Hash(password)
	if err != nil {
		return "", err
	}

	// Encrypt xpriv with hashed password
	encryptedXpriv, err := encryption.Encrypt(hashedPassword, xpriv.String())

	if err != nil {
		return "", err
	}

	user := &User{
		Email:     email,
		Xpriv:     encryptedXpriv,
		CreatedAt: time.Now(),
	}

	err = s.InsertUser(user)
	if err == nil {
		err = s.BuxClient.RegisterXpub(xpriv)
		if err != nil {
			return "", err
		}
	}

	return mnemonic, err
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
