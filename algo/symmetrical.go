package algo

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

type (
	Symmetrical struct {
	}
)

func NewSymmetrical() *Symmetrical {
	return &Symmetrical{}
}

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

func MakeSimpleHash(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}
