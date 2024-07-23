package security_test

import (
	"SecureChat/internal/security"
	"testing"
)

func TestComputeSharedSecret(t *testing.T) {
	//arrange
	alicePrivate, alicePublic, err := security.GenerateKeyPair()
	if err != nil {
		panic(err)
	}
	bobPrivate, bobPublic, err := security.GenerateKeyPair()
	if err != nil {
		panic(err)
	}

	// act
	aliceSharedSecret, err := security.ComputeSharedSecret(alicePrivate, bobPublic)
	if err != nil {
		panic(err)
	}
	bobSharedSecret, err := security.ComputeSharedSecret(bobPrivate, alicePublic)
	if err != nil {
		panic(err)
	}

	//assert
	if aliceSharedSecret != bobSharedSecret {
		t.Errorf("different shared secrets alice: %s bob: %s", aliceSharedSecret[:32], bobSharedSecret[:32])
	}
}
