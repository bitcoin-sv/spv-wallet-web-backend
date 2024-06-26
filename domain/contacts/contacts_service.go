package contacts

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type ContactsService struct {
	adminWalletClient   users.AdminWalletClient
	walletClientFactory users.WalletClientFactory
	log                 *zerolog.Logger
}

func NewContactsService(adminWalletClient users.AdminWalletClient, walletClientFactory users.WalletClientFactory, log *zerolog.Logger) *ContactsService {
	transactionServiceLogger := log.With().Str("service", "contacts-service").Logger()
	return &ContactsService{
		adminWalletClient:   adminWalletClient,
		walletClientFactory: walletClientFactory,
		log:                 &transactionServiceLogger,
	}
}

func (s *ContactsService) UpsertContact(ctx context.Context, accessKey, paymail, fullName, requesterPaymail string, metadata map[string]any) (*models.Contact, error) {
	userWalletClient, err := s.walletClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return nil, err
	}

	return userWalletClient.UpsertContact(ctx, paymail, fullName, requesterPaymail, metadata)
}

func (s *ContactsService) AcceptContact(ctx context.Context, accessKey, paymail string) error {
	userWalletClient, err := s.walletClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return err
	}

	return userWalletClient.AcceptContact(ctx, paymail)
}

func (s *ContactsService) RejectContact(ctx context.Context, accessKey, paymail string) error {
	userWalletClient, err := s.walletClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return err
	}

	return userWalletClient.RejectContact(ctx, paymail)
}

func (s *ContactsService) ConfirmContact(ctx context.Context, xPriv string, contact *models.Contact, passcode, requesterPaymail string) error {
	userWalletClient, err := s.walletClientFactory.CreateWithXpriv(xPriv)
	if err != nil {
		return err
	}

	return userWalletClient.ConfirmContact(ctx, contact, passcode, requesterPaymail, getConfPeriod(), getConfDigits())
}

func (s *ContactsService) GetContacts(ctx context.Context, accessKey string, conditions *filter.ContactFilter, metadata map[string]any, queryParams *filter.QueryParams) (*models.SearchContactsResponse, error) {
	userWalletClient, err := s.walletClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return nil, err
	}

	return userWalletClient.GetContacts(ctx, conditions, metadata, queryParams)
}

func (s *ContactsService) GenerateTotpForContact(ctx context.Context, xPriv string, contact *models.Contact) (string, error) {
	userWalletClient, err := s.walletClientFactory.CreateWithXpriv(xPriv) //xPriv instead of accessKey because it is necessary to calculate the shared secret
	if err != nil {
		return "", err
	}

	return userWalletClient.GenerateTotpForContact(contact, getConfPeriod(), getConfDigits())
}

func getConfPeriod() uint {
	return viper.GetUint(config.EnvContactsPasscodePeriod)
}
func getConfDigits() uint {
	return viper.GetUint(config.EnvContactsPasscodeDigits)
}
