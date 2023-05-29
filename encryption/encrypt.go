package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/xdg-go/pbkdf2"
)

func deriveKey(passphrase string, salt []byte) ([]byte, []byte) {
	return pbkdf2.Key([]byte(passphrase), salt, 1000, 32, sha256.New), salt
}

// Encrypt encrypts the plaintext using AES-GCM.
func Encrypt(passphrase, plaintext string) (string, error) {
	key, salt := deriveKey(passphrase, nil)
	iv := make([]byte, 12)
	_, err := rand.Read(iv)
	if err != nil {
		panic(err)
	}
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	data := aesgcm.Seal(nil, iv, []byte(plaintext), nil)
	return hex.EncodeToString(salt) + "-" + hex.EncodeToString(iv) + "-" + hex.EncodeToString(data), nil
}

// Decrypt decrypts the ciphertext using AES-GCM.
func Decrypt(passphrase, ciphertext string) string {
	arr := strings.Split(ciphertext, "-")
	salt, _ := hex.DecodeString(arr[0])
	iv, _ := hex.DecodeString(arr[1])
	data, _ := hex.DecodeString(arr[2])
	key, _ := deriveKey(passphrase, salt)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	data, _ = aesgcm.Open(nil, iv, data, nil)
	return string(data)
}
