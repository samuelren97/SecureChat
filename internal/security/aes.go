package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

var (
	ErrIncorrectKeyLength error = errors.New("incorrect key length")
)

func EncryptAESGCM(plainText, key string) (string, error) {
	plainTextBytes := []byte(plainText)
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return "", ErrIncorrectKeyLength
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, plainTextBytes, nil)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func DecryptAESGCM(cipherText, key string) (string, error) {
	cipherTextBytes, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	keyBytes := []byte(key)
	if len(key) != 32 {
		return "", ErrIncorrectKeyLength
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, cipherTextBytes := cipherTextBytes[:nonceSize], cipherTextBytes[nonceSize:]

	plainTextBytes, err := gcm.Open(nil, nonce, cipherTextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plainTextBytes), nil
}

func GenerateKey(length int) string {
	byteSlice := make([]byte, length)
	rand.Read(byteSlice)

	return hex.EncodeToString(byteSlice)
}
