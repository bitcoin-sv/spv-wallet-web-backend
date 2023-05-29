package hash

import (
	"fmt"
	"hash"
)

// SHA256Hasher is a SHA256 password hasher.
type SHA256Hasher struct {
	hash hash.Hash
	salt string
}

// NewSHA256Hasher creates a new SHA256Hasher.
func NewSHA256Hasher(salt string, hash hash.Hash) *SHA256Hasher {
	return &SHA256Hasher{
		salt: salt,
		hash: hash,
	}
}

// Hash hashes the data using SHA256.
func (h *SHA256Hasher) Hash(data string) (string, error) {
	if _, err := h.hash.Write([]byte(data)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.hash.Sum([]byte(h.salt))), nil
}
