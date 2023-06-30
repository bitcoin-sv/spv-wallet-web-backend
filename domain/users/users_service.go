package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"bux-wallet/config"
	"bux-wallet/encryption"
	"bux-wallet/logging"

	"github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/bip39"
	"github.com/libsv/go-bk/chaincfg"
	"github.com/spf13/viper"
)

var WrongUserPassword = errors.New("Wrong Password")

// UserService represents User service and provide access to repository.
type UserService struct {
	repo             UsersRepository
	buxClient        AdmBuxClient
	buxClientFactory BuxClientFactory
	log              logging.Logger
}

// NewUserService creates UserService instance.
func NewUserService(repo UsersRepository, adminBuxClient AdmBuxClient, bf BuxClientFactory, lf logging.LoggerFactory) *UserService {
	// Create service.
	s := &UserService{
		repo:             repo,
		buxClient:        adminBuxClient,
		buxClientFactory: bf,
		log:              lf.NewLogger("user-service"),
	}

	return s
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
	xpub, err := s.buxClient.RegisterXpub(xpriv)
	if err != nil {
		return nil, fmt.Errorf("error registering xpub in BUX: %s", err.Error())
	}

	// Get username from email which will be used as paymail alias.
	username, _ := splitEmail(email)

	// Register paymail in BUX.
	paymail, err := s.buxClient.RegisterPaymail(username, xpub)
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
	buxClient, err := s.buxClientFactory.CreateWithXpriv(decryptedXpriv)
	if err != nil {
		return nil, err
	}

	// Create access key.
	accessKey, err := buxClient.CreateAccessKey()
	if err != nil {
		return nil, err
	}

	xpub, err := buxClient.GetXPub()
	if err != nil {
		return nil, err
	}

	balance, err := calculateBalance(xpub.GetCurrentBalance())
	if err != nil {
		return nil, err
	}

	signInUser := &AuthenticatedUser{
		User: user,
		AccessKey: AccessKey{
			Id:  accessKey.GetAccessKeyId(),
			Key: accessKey.GetAccessKey(),
		},
		Balance: *balance,
	}

	return signInUser, nil
}

// SignOutUser signs out user by revoking access key. (Not possible at the moment, method is just a mock.)
func (s *UserService) SignOutUser(accessKeyId, accessKey string) error {

	/// Right now we cannot revoke access key without Bux Client authentication with XPriv, which is impossible here.

	// buxClient, err := s.buxClientFactory.CreateWithAccessKey(accessKey)
	// if err != nil {
	// 	return err
	// }

	// _, err = buxClient.RevokeAccessKey(accessKeyId)
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

// GetUserBalance returns user balance. Bux client is created with access key.
func (s *UserService) GetUserBalance(accessKey string) (*Balance, error) {
	// Create BUX client with access key.
	buxClient, err := s.buxClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return nil, err
	}

	// Get xpub.
	xpub, err := buxClient.GetXPub()
	if err != nil {
		return nil, err
	}

	// Calculate balance.
	balance, err := calculateBalance(xpub.GetCurrentBalance())
	if err != nil {
		return nil, err
	}

	return balance, nil
}

// GetUserXpriv gets user by id and decrypt xpriv.
func (s *UserService) GetUserXpriv(userId int, password string) (string, error) {
	user, err := s.repo.GetUserById(context.Background(), userId)
	if err != nil {
		return "", err
	}

	// Decrypt xpriv.
	decryptedXpriv, err := decryptXpriv(password, user.Xpriv)
	if err != nil {
		return "", err
	}

	return decryptedXpriv, nil
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
	fmt.Println("xpriv", xpriv)
	fmt.Println("password", password)
	// Create hash from password
	hashedPassword, err := encryption.Hash(password)
	if err != nil {
		return "", err
	}

	fmt.Println("hashedPassword", hashedPassword)

	// Encrypt xpriv with hashed password
	encryptedXpriv, err := encryption.Encrypt(hashedPassword, xpriv)
	if err != nil {
		return "", err
	}

	fmt.Println("encryptedXpriv", encryptedXpriv)

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
	if xpriv == "" {
		return "", WrongUserPassword
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

func calculateBalance(satoshis uint64) (*Balance, error) {
	// Create request.
	exchangeRateUrl := viper.GetString(config.EnvEndpointsExchangeRate)
	req, error := http.NewRequestWithContext(context.Background(), http.MethodGet, exchangeRateUrl, nil)
	if error != nil {
		return nil, fmt.Errorf("error during creating exchange rate request: %s", error.Error())
	}

	// Send request.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during getting exchange rate: %s", err.Error())
	}
	defer res.Body.Close() // nolint: all

	// Parse response body.
	var exchangeRate ExchangeRate
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error during reading response body: %s", err.Error())
	}

	err = json.Unmarshal(bodyBytes, &exchangeRate)
	if err != nil {
		return nil, fmt.Errorf("error during unmarshaling response body: %s", err.Error())
	}

	// Calculate balance.
	balanceBSV := float64(satoshis) / 100000000
	balanceUSD := balanceBSV * exchangeRate.Rate

	balance := &Balance{
		Bsv:      balanceBSV,
		Usd:      balanceUSD,
		Satoshis: satoshis,
	}

	return balance, nil
}
