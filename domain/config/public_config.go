package config

// PublicConfig represents a config that is exposed to the public.
type PublicConfig struct {
	PaymilDomain         string          `json:"paymail_domain"`
	ExperimentalFeatures map[string]bool `json:"experimental_features"`
}
