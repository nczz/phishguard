package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

// Encrypt encrypts plaintext using AES-256-GCM. Key must be 64-char hex (32 bytes).
// Returns nonce+ciphertext as raw bytes (suitable for VARBINARY storage).
func Encrypt(hexKey, plaintext string) ([]byte, error) {
	if plaintext == "" {
		return nil, nil
	}
	k, err := hex.DecodeString(hexKey)
	if err != nil || len(k) != 32 {
		return nil, errors.New("ENCRYPT_KEY must be 64-char hex (32 bytes)")
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, []byte(plaintext), nil), nil
}

// Decrypt decrypts AES-256-GCM ciphertext (nonce prepended).
// Falls back to treating data as plaintext if decryption fails (backward compat with pre-encryption data).
func Decrypt(hexKey string, ciphertext []byte) (string, error) {
	if len(ciphertext) == 0 {
		return "", nil
	}
	k, err := hex.DecodeString(hexKey)
	if err != nil || len(k) != 32 {
		return "", errors.New("ENCRYPT_KEY must be 64-char hex (32 bytes)")
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		// Too short to be encrypted — treat as legacy plaintext
		return string(ciphertext), nil
	}
	plaintext, err := gcm.Open(nil, ciphertext[:nonceSize], ciphertext[nonceSize:], nil)
	if err != nil {
		// Decryption failed — assume legacy plaintext data
		return string(ciphertext), nil
	}
	return string(plaintext), nil
}
