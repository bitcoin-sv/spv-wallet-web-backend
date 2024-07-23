package contacts

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet-web-backend/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
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
	userWalletClient := s.walletClientFactory.CreateWithAccessKey(accessKey)

	contact, err := userWalletClient.UpsertContact(ctx, paymail, fullName, requesterPaymail, metadata)
	if err != nil {
		s.log.Debug().Msgf("Error during upsert contact: %s", err.Error())
		return nil, spverrors.ErrUpsertContact
	}
	return contact, nil
}

// AcceptContact accepts a contact invitation
func (s *Service) AcceptContact(ctx context.Context, accessKey, paymail string) error {
	userWalletClient := s.walletClientFactory.CreateWithAccessKey(accessKey)

	err := userWalletClient.AcceptContact(ctx, paymail)
	if err != nil {
		s.log.Debug().Msgf("Error during accepting contact: %s", err.Error())
		return spverrors.ErrAcceptContact
	}
	return nil
}

// RejectContact rejects a contact invitation
func (s *Service) RejectContact(ctx context.Context, accessKey, paymail string) error {
	userWalletClient := s.walletClientFactory.CreateWithAccessKey(accessKey)

	err := userWalletClient.RejectContact(ctx, paymail)
	if err != nil {
		s.log.Debug().Msgf("Error during rejecting contact: %s", err.Error())
		return spverrors.ErrRejectContact
	}
	return nil
}

// ConfirmContact confirms a contact
func (s *Service) ConfirmContact(ctx context.Context, xPriv string, contact *models.Contact, passcode, requesterPaymail string) error {
	userWalletClient := s.walletClientFactory.CreateWithXpriv(xPriv)

	err := userWalletClient.ConfirmContact(ctx, contact, passcode, requesterPaymail, getConfPeriod(), getConfDigits())
	if err != nil {
		s.log.Debug().Msgf("Error during confirming contact: %s", err.Error())
		return spverrors.ErrConfirmContact
	}
	return nil
}

// GetContacts retrieves contacts for the user
func (s *Service) GetContacts(ctx context.Context, accessKey string, conditions *filter.ContactFilter, metadata map[string]any, queryParams *filter.QueryParams) (*models.SearchContactsResponse, error) {
	userWalletClient := s.walletClientFactory.CreateWithAccessKey(accessKey)

	resp, err := userWalletClient.GetContacts(ctx, conditions, metadata, queryParams)
	if err != nil {
		s.log.Debug().Msgf("Error during getting contacts: %s", err.Error())
		return nil, spverrors.ErrGetContacts
	}
	return resp, nil
}

// GenerateTotpForContact generates a TOTP for a contact
func (s *Service) GenerateTotpForContact(_ context.Context, xPriv string, contact *models.Contact) (string, error) {
	userWalletClient := s.walletClientFactory.CreateWithXpriv(xPriv) // xPriv instead of accessKey because it is necessary to calculate the shared secret

	totp, err := userWalletClient.GenerateTotpForContact(contact, getConfPeriod(), getConfDigits())
	if err != nil {
		s.log.Debug().Msgf("Error during generating TOTP for contact: %s", err.Error())
		return "", spverrors.ErrGenerateTotpForContact
	}
	return totp, nil
}

func getConfPeriod() uint {
	return viper.GetUint(config.EnvContactsPasscodePeriod)
}
func getConfDigits() uint {
	return viper.GetUint(config.EnvContactsPasscodeDigits)
}
