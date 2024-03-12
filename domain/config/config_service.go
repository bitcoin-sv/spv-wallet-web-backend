package config

import (
	"slices"
	"sync"

	backendconfig "github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// ConfigService is a service for fetching and caching SharedConfig from the spv-wallet and providing PublicConfig.
type ConfigService struct {
	adminWalletClient users.AdminWalletClient
	log               *zerolog.Logger

	sharedConfig *models.SharedConfig
	publicConfig *PublicConfig
	mutex        sync.Mutex
}

// NewConfigService creates a new ConfigService.
func NewConfigService(adminWalletClient users.AdminWalletClient, log *zerolog.Logger) *ConfigService {
	return &ConfigService{
		adminWalletClient: adminWalletClient,
		log:               log,
		sharedConfig:      nil,
		publicConfig:      nil,
	}
}

// GetSharedConfig returns shared config.
// If shared config is not cached, it will be fetched from the spv-wallet.
// SharedConfig should not be exposed to the public - use PublicConfig instead.
func (s *ConfigService) GetSharedConfig() *models.SharedConfig {
	if s.sharedConfig != nil {
		return s.sharedConfig
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	model, err := s.adminWalletClient.GetSharedConfig()
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get shared config")
		return nil
	}
	s.sharedConfig = model
	return s.sharedConfig
}

// GetPublicConfig returns public config.
func (s *ConfigService) GetPublicConfig() *PublicConfig {
	if s.publicConfig != nil {
		return s.publicConfig
	}
	shared := s.GetSharedConfig()
	if shared == nil {
		return nil
	}

	s.publicConfig = s.makePublicConfig(shared)
	return s.publicConfig
}

func (s *ConfigService) makePublicConfig(shared *models.SharedConfig) *PublicConfig {
	configuredPaymailDomain := viper.GetString(backendconfig.EnvPaymailDomain)
	if !slices.Contains(shared.PaymilDomains, configuredPaymailDomain) {
		s.log.Warn().Str("configuredPaymailDomain", configuredPaymailDomain).Msg("Configured paymail domain is not in the list of paymail domains from SPV Wallet")
	}

	return &PublicConfig{
		PaymilDomain:         configuredPaymailDomain,
		ExperimentalFeatures: shared.ExperimentalFeatures,
	}
}
