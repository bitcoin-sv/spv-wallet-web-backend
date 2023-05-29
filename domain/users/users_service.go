package users

import (
	"bux-wallet/data/users"
	"bux-wallet/hash"
	bux_client "bux-wallet/transports/bux/client"
	"context"
	"time"

	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/bip39"
	"github.com/libsv/go-bk/chaincfg"
)

// UserService represents User service and provide access to repository.
type UserService struct {
	repo      UsersRepository
	BuxClient *bux_client.BClient
	Hasher    *hash.SHA256Hasher
}

// NewUserService creates UserService instance.
func NewUserService(repo *users.UsersRepository, buxClient *bux_client.BClient, hasher *hash.SHA256Hasher) *UserService {
	return &UserService{
		repo:      repo,
		BuxClient: buxClient,
		Hasher:    hasher,
	}
}

// InsertUser inserts user to database.
func (s *UserService) InsertUser(user *User) error {
	err := s.repo.InsertUser(context.Background(), user.toUserDto())
	return err
}

// CreateNewUser creates new user.
func (s *UserService) CreateNewUser(email, password string) (*User, error) {
	// Generate mnemonic and seed
	mnemonic, seed, err := generateMnemonic()
	if err != nil {
		return nil, err
	}

	xpriv, err := generateXpriv(seed)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:     email,
		Password:  password,
		Mnemonic:  mnemonic,
		Xpriv:     xpriv.String(),
		CreatedAt: time.Now(),
	}

	// Hash user data
	err = s.hashUser(user)
	if err != nil {
		return nil, err
	}

	err = s.InsertUser(user)
	if err == nil {
		s.BuxClient.RegisterXpub(xpriv)

	}

	// Restore uncrypted mnemonic to show it to user.
	user.Mnemonic = mnemonic
	return user, err
}

func (s *UserService) hashUser(user *User) error {
	hashedPassword, err := s.Hasher.Hash(user.Password)
	if err != nil {
		return err
	}
	hashedMnemonic, err := s.Hasher.Hash(user.Mnemonic)
	if err != nil {
		return err
	}
	hashedXpriv, err := s.Hasher.Hash(user.Xpriv)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.Mnemonic = hashedMnemonic
	user.Xpriv = hashedXpriv

	return nil
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
