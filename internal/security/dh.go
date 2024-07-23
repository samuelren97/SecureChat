package security

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

var curve = ecdh.P256()

func GenerateKeyPair() (string, string, error) {
	priv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	return base64.URLEncoding.EncodeToString(priv.Bytes()),
		base64.URLEncoding.EncodeToString(priv.PublicKey().Bytes()),
		nil
}

func ComputeSharedSecret(priv, peerPub string) (string, error) {
	pubBytes, err := base64.URLEncoding.DecodeString(peerPub)
	if err != nil {
		return "", err
	}
	pubKey, err := curve.NewPublicKey(pubBytes)
	if err != nil {
		return "", err
	}

	privBytes, err := base64.URLEncoding.DecodeString(priv)
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

	return base64.URLEncoding.EncodeToString(sharedSecret), nil
}
