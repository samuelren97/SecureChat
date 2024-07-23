package security_test

import (
	"SecureChat/internal/security"
	"testing"
)

func TestEncryptAESGCM_DifferentResultSameKey(t *testing.T) {
	//arrange
	plainText := "hello goland!"
	key := security.GenerateKey(16)

	//act
	cipherText1, err := security.EncryptAESGCM(plainText, key)
	if err != nil {
		panic(err)
	}
	cipherText2, err := security.EncryptAESGCM(plainText, key)
	if err != nil {
		panic(err)
	}

	//assert
	if cipherText1 == cipherText2 {
		t.Error("the ciphertexts should not be the same")
	}
}

func TestDecryptAESGCM_OneCipherText(t *testing.T) {
	//arrange
	expectedPlainText := "hello goland!"
	key := security.GenerateKey(16)
	cipherText, err := security.EncryptAESGCM(expectedPlainText, key)
	if err != nil {
		panic(err)
	}

	//act
	result, err := security.DecryptAESGCM(cipherText, key)
	if err != nil {
		panic(err)
	}

	//assert
	if result != expectedPlainText {
		t.Errorf("wanted: %s got: %s", expectedPlainText, result)
	}
}

func TestDecryptAESGCM_TwoCipherTextsSameSourceSameKey(t *testing.T) {
	//arrange
	expectedPlainText := "hello goland!"
	key := security.GenerateKey(16)
	cipherText1, err := security.EncryptAESGCM(expectedPlainText, key)
	if err != nil {
		panic(err)
	}
	cipherText2, err := security.EncryptAESGCM(expectedPlainText, key)
	if err != nil {
		panic(err)
	}

	//act
	result1, err := security.DecryptAESGCM(cipherText1, key)
	if err != nil {
		panic(err)
	}

	result2, err := security.DecryptAESGCM(cipherText2, key)
	if err != nil {
		panic(err)
	}

	//assert
	if result1 != result2 {
		t.Errorf("different results => result1: %s result2: %s", result1, result2)
	}
}
