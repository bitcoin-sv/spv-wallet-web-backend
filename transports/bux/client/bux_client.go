package buxclient

import (
	"bux-wallet/logging"
)

// DefaultBuxClientFactory creates an instance of Bux Client Factory.
func DefaultBuxClientFactory(lf logging.LoggerFactory) BuxClientFactory {
	return NewBuxClientFactory(lf)
}
