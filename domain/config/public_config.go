package config

import "github.com/bitcoin-sv/spv-wallet/models"

type PublicConfig struct {
	PaymilDomain         string                    `json:"paymail_domain"`
	ExperimentalFeatures models.ExperimentalConfig `json:"experimental_features"`
}
