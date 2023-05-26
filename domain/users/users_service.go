package users

import (
	"bux-wallet/data/users"
	"context"
	"time"

	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/bip39"
	"github.com/libsv/go-bk/chaincfg"
)

// UserService represents User service and provide access to repository.
type UserService struct {
	repo UsersRepository
}

// NewUserService creates UserService instance.
func NewUserService(repo *users.UsersRepository) *UserService {
	return &UserService{
		repo: repo,
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
		Xpriv:     xpriv,
		CreatedAt: time.Now(),
	}
	err = s.InsertUser(user)
	return user, err
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
func generateXpriv(seed []byte) (string, error) {
	hdXpriv, err := bip32.NewMaster(seed, &chaincfg.MainNet)
	if err != nil {
		return "", err
	}
	return hdXpriv.String(), nil
}
