package security

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"errors"
)

var (
	ErrIncorrectKeyLength error = errors.New("incorrect key length")
)

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func Encrypt(message string, key string) (string, error) {
	if len(key) != 32 {
		return "", ErrIncorrectKeyLength
	}

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Pad the message to be a multiple of the block size
	paddedMessage := pkcs7Pad([]byte(message), aes.BlockSize)

	out := make([]byte, len(paddedMessage))
	c.Encrypt(out, paddedMessage)

	return hex.EncodeToString(out), nil
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("decryption error: padding size is zero")
	}

	padding := int(data[length-1])
	if padding > length || padding > aes.BlockSize {
		return nil, errors.New("decryption error: invalid padding size")
	}

	for i := 0; i < padding; i++ {
		if data[length-1-i] != byte(padding) {
			return nil, errors.New("decryption error: invalid padding")
		}
	}

	return data[:length-padding], nil
}

func Decrypt(message string, key string) (string, error) {
	if len(key) != 32 {
		return "", ErrIncorrectKeyLength
	}

	cipherText, err := hex.DecodeString(message)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	pt := make([]byte, len(cipherText))
	c.Decrypt(pt, cipherText)

	// Unpad the decrypted message
	unpaddedMessage, err := pkcs7Unpad(pt)
	if err != nil {
		return "", err
	}

	return string(unpaddedMessage), nil
}
