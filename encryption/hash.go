package encryption

import (
	"crypto/sha256"
	"fmt"
	"web-backend/config"

	"github.com/spf13/viper"
)

// Hash hashes the data using SHA256.
func Hash(data string) (string, error) {
	hash := sha256.New()
	salt := viper.GetString(config.EnvHashSalt)

	if _, err := hash.Write([]byte(data)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(salt))), nil
}
