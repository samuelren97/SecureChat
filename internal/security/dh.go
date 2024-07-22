package security

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

var curve = ecdh.P256()

func GenerateKeyPair() (string, string, error) {
	priv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	return hex.EncodeToString(priv.Bytes()), hex.EncodeToString(priv.PublicKey().Bytes()), nil
}

func ComputeSharedSecret(priv, peerPub string) (string, error) {
	pubBytes, err := hex.DecodeString(peerPub)
	if err != nil {
		return "", err
	}
	pubKey, err := curve.NewPublicKey(pubBytes)
	if err != nil {
		return "", err
	}

	privBytes, err := hex.DecodeString(priv)
	if err != nil {
		return "", err
	}

	privKey, err := curve.NewPrivateKey(privBytes)
	if err != nil {
		return "", err
	}

	sharedSecret, err := privKey.ECDH(pubKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sharedSecret), nil
}
