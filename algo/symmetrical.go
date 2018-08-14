package algo

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

type (
	// Symmetrical represents a collections of encryption and decryption symmetrical algorithms
	// A symmetrical algorithm is the one that uses a single key to convert data
	Symmetrical struct {
	}
)

// NewSymmetrical returns a reference to the type and access to its funcionalities
func NewSymmetrical() *Symmetrical {
	return &Symmetrical{}
}

// Encrypt uses a passphrase to encrypt data using gcm algorithm
func (s *Symmetrical) Encrypt(data []byte, passphrase string) ([]byte, error) {

	block, _ := aes.NewCipher(MakeSimpleHash(passphrase))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt uses the encription passphrase to decrypt data to its original state
func (s *Symmetrical) Decrypt(data []byte, passphrase string) ([]byte, error) {

	key := MakeSimpleHash(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// MakeSimpleHash returns a bytes array with a sha256 hash encryption
func MakeSimpleHash(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}
