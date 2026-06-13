package encryptservice

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"ghivarra/afterglow/environment"
	"io"
)

func Encrypt(rawData string) (string, error) {
	var key = []byte(environment.ENCRYPTION_KEY)

	if len(key) != 32 {
		return "", errors.New("key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nil, nonce, []byte(rawData), nil)

	// store nonce + ciphertext together
	finalData := append(nonce, cipherText...)

	return base64.StdEncoding.EncodeToString(finalData), nil
}

func Decrypt(encryptedData string) (string, error) {
	var key = []byte(environment.ENCRYPTION_KEY)

	if len(key) != 32 {
		return "", errors.New("key must be 32 bytes for AES-256")
	}

	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("invalid encrypted text")
	}

	nonce := data[:nonceSize]
	cipherText := data[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
