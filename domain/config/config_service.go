package config

import (
	"slices"
	"sync"
	"time"

	backendconfig "github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

const cacheTTL = 10 * time.Minute

// ConfigService is a service for fetching and caching SharedConfig from the spv-wallet and providing PublicConfig.
type ConfigService struct {
	adminWalletClient users.AdminWalletClient
	log               *zerolog.Logger

	sharedConfig *models.SharedConfig
	publicConfig *PublicConfig
	mutex        sync.Mutex
	lastFetch    time.Time
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
	s.makeConfigs()
	return s.sharedConfig
}

// GetPublicConfig returns public config.
func (s *ConfigService) GetPublicConfig() *PublicConfig {
	s.makeConfigs()
	return s.publicConfig
}

func (s *ConfigService) makeConfigs() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.sharedConfig != nil && time.Since(s.lastFetch) < cacheTTL {
		return
	}
	s.lastFetch = time.Now()

	shared, err := s.adminWalletClient.GetSharedConfig()
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get shared config")
		return
	}

	s.sharedConfig = shared
	s.publicConfig = s.makePublicConfig(shared)
}

func (s *ConfigService) makePublicConfig(shared *models.SharedConfig) *PublicConfig {
	configuredPaymailDomain := viper.GetString(backendconfig.EnvPaymailDomain)
	if !slices.Contains(shared.PaymailDomains, configuredPaymailDomain) {
		s.log.Warn().Str("configuredPaymailDomain", configuredPaymailDomain).Msg("Configured paymail domain is not in the list of paymail domains from SPV Wallet")
	}

	return &PublicConfig{
		PaymilDomain:         configuredPaymailDomain,
		ExperimentalFeatures: shared.ExperimentalFeatures,
	}
}
