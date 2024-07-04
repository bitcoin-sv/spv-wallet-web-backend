package encryption

import (
	"crypto/sha256"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Hash hashes the data using SHA256.
func Hash(data string) (string, error) {
	hash := sha256.New()
	salt := viper.GetString(config.EnvHashSalt)

	if _, err := hash.Write([]byte(data)); err != nil {
		return "", errors.Wrap(err, "internal error")
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(salt))), nil
}
