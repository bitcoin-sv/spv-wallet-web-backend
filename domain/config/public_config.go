package config

// PublicConfig represents a config that is exposed to the public.
type PublicConfig struct {
	PaymilDomain         string             `json:"paymail_domain"`
	ExperimentalFeatures ExperimentalConfig `json:"experimental_features"`
}

// ExperimentalConfig represents a feature flag config.
type ExperimentalConfig struct {
	PikeEnabled bool `json:"pike_enabled"`
}
