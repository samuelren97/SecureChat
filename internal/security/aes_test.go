package security_test

import (
	"SecureChat/internal/security"
	"errors"
	"testing"
)

func TestEncryptAESGCM_DifferentResultSameKey(t *testing.T) {
	//arrange
	plainText := "hello goland!"
	key := security.GenerateSecureKey(32)

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

func TestEncryptAESGCM_IncorrectKeyLength_Err(t *testing.T) {
	//arrange
	plainText := "hello goland!"
	key := security.GenerateSecureKey(31)

	//act
	_, err := security.EncryptAESGCM(plainText, key)

	//assert
	if err == nil {
		t.Error("Expected error but it did not happen")
	} else if !errors.Is(err, security.ErrIncorrectKeyLength) {
		t.Errorf("An error was returned but is not of incorrect key length: %v", err.Error())
	}
}

func TestDecryptAESGCM_OneCipherText(t *testing.T) {
	//arrange
	expectedPlainText := "hello goland!"
	key := security.GenerateSecureKey(32)
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
	key := security.GenerateSecureKey(32)
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

func TestDecryptAESGCM_IncorrectKeyLength_Err(t *testing.T) {
	//arrange
	cipherText := "aGVsbG8gZ29sYW5kIQ=="
	key := security.GenerateSecureKey(31)

	//act
	_, err := security.DecryptAESGCM(cipherText, key)

	//assert
	if err == nil {
		t.Error("Expected error but it did not happen")
	} else if !errors.Is(err, security.ErrIncorrectKeyLength) {
		t.Errorf("An error was returned but is not of incorrect key length: %v", err.Error())
	}
}
