package encryption_test

import (
	"bux-wallet/encryption"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEncryptionDecryption tests if we use same algorithm for encryption and decryption data
func TestEncryptionDecryption(t *testing.T) {
	cases := []struct {
		name                           string
		plaintext                      string
		passphraseForEncryption        string
		passphraseForDecryption        string
		expectedDecodedTextSameAsInput bool
	}{
		{
			name:                           "Used valid passphrase",
			plaintext:                      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque consequat justo sed viverra laoreet. Cras pulvinar orci turpis, vel pulvinar tellus vestibulum eget. Duis dictum purus eget ante aliquet mollis. Duis sit amet arcu urna.",
			passphraseForEncryption:        "4chain",
			passphraseForDecryption:        "4chain",
			expectedDecodedTextSameAsInput: true,
		},
		{
			name:                           "Used invalid passphrase",
			plaintext:                      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque consequat justo sed viverra laoreet. Cras pulvinar orci turpis, vel pulvinar tellus vestibulum eget. Duis dictum purus eget ante aliquet mollis. Duis sit amet arcu urna.",
			passphraseForEncryption:        "4chain",
			passphraseForDecryption:        "otherword",
			expectedDecodedTextSameAsInput: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			ciphertext, _ := encryption.Encrypt(tc.passphraseForEncryption, tc.plaintext)
			decodedtext := encryption.Decrypt(tc.passphraseForDecryption, ciphertext)

			// Assert
			if tc.expectedDecodedTextSameAsInput {
				assert.Equal(t, tc.plaintext, decodedtext)
			} else {
				assert.NotEqual(t, tc.plaintext, decodedtext)
			}

		})
	}
}

// TestHash tests if SHA256 is used correctly
func TestHash(t *testing.T) {
	tc := struct {
		name         string
		input        string
		expectedHash string
	}{
		name:         "Test SHA256",
		input:        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque consequat justo sed viverra laoreet. Cras pulvinar orci turpis, vel pulvinar tellus vestibulum eget. Duis dictum purus eget ante aliquet mollis. Duis sit amet arcu urna.",
		expectedHash: "4c87fbdd08e7e87f71b7f9ca22f0cf60206f604d00acced0656bb28e7d640482", // generated with online SHA256 hasher
	}

	t.Run(tc.name, func(t *testing.T) {
		hash, _ := encryption.Hash(tc.input)

		assert.Equal(t, tc.expectedHash, hash)
	})
}
