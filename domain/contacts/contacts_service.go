package contacts

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Service is the service that manages contacts
type Service struct {
	adminWalletClient   users.AdminWalletClient
	walletClientFactory users.WalletClientFactory
	log                 *zerolog.Logger
}

// NewContactsService creates a new instance of the contact.Service
func NewContactsService(adminWalletClient users.AdminWalletClient, walletClientFactory users.WalletClientFactory, log *zerolog.Logger) *Service {
	transactionServiceLogger := log.With().Str("service", "contacts-service").Logger()
	return &Service{
		adminWalletClient:   adminWalletClient,
		walletClientFactory: walletClientFactory,
		log:                 &transactionServiceLogger,
	}
}

// UpsertContact creates or updates a contact
func (s *Service) UpsertContact(ctx context.Context, accessKey, paymail, fullName, requesterPaymail string, metadata map[string]any) (*models.Contact, error) {
	userWalletClient, err := s.walletClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return nil, errors.Wrap(err, "internal error")
	}

	contact, err := userWalletClient.UpsertContact(ctx, paymail, fullName, requesterPaymail, metadata)
	if err != nil {
		return nil, errors.Wrap(err, "upsert contact error")
	}
	return contact, nil
}

// AcceptContact accepts a contact invitation
func (s *Service) AcceptContact(ctx context.Context, accessKey, paymail string) error {
	userWalletClient, err := s.walletClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return errors.Wrap(err, "internal error")
	}

	return errors.Wrap(userWalletClient.AcceptContact(ctx, paymail), "accept contact error")
}

// RejectContact rejects a contact invitation
func (s *Service) RejectContact(ctx context.Context, accessKey, paymail string) error {
	userWalletClient, err := s.walletClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return errors.Wrap(err, "internal error")
	}

	return errors.Wrap(userWalletClient.RejectContact(ctx, paymail), "reject contact error")
}

// ConfirmContact confirms a contact
func (s *Service) ConfirmContact(ctx context.Context, xPriv string, contact *models.Contact, passcode, requesterPaymail string) error {
	userWalletClient, err := s.walletClientFactory.CreateWithXpriv(xPriv)
	if err != nil {
		return errors.Wrap(err, "internal error")
	}

	return errors.Wrap(userWalletClient.ConfirmContact(ctx, contact, passcode, requesterPaymail, getConfPeriod(), getConfDigits()), "confirm contact error")
}

// GetContacts retrieves contacts for the user
func (s *Service) GetContacts(ctx context.Context, accessKey string, conditions *filter.ContactFilter, metadata map[string]any, queryParams *filter.QueryParams) (*models.SearchContactsResponse, error) {
	userWalletClient, err := s.walletClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return nil, errors.Wrap(err, "internal error")
	}

	resp, err := userWalletClient.GetContacts(ctx, conditions, metadata, queryParams)
	if err != nil {
		return nil, errors.Wrap(err, "get contacts error")
	}
	return resp, nil
}

// GenerateTotpForContact generates a TOTP for a contact
func (s *Service) GenerateTotpForContact(_ context.Context, xPriv string, contact *models.Contact) (string, error) {
	userWalletClient, err := s.walletClientFactory.CreateWithXpriv(xPriv) // xPriv instead of accessKey because it is necessary to calculate the shared secret
	if err != nil {
		return "", errors.Wrap(err, "internal error")
	}

	totp, err := userWalletClient.GenerateTotpForContact(contact, getConfPeriod(), getConfDigits())
	return totp, errors.Wrap(err, "spv wallet error")
}

func getConfPeriod() uint {
	return viper.GetUint(config.EnvContactsPasscodePeriod)
}
func getConfDigits() uint {
	return viper.GetUint(config.EnvContactsPasscodeDigits)
}
