package security_test

import (
	"NewNASAPI/internal/security"
	"testing"
)

func TestCompareHashes(t *testing.T) {
	// Arrange
	input := "password"

	hash := &security.Hash{
		Hash: "a8a4fedfb336b72dd0d7853a65a6ec8c0a2b20f2dce22a2cdcf0bd59ecc37ee8",
		Salt: "f42a97f5564970bd013b8261eb41843b",
	}

	// Act
	isSameHash := security.CompareHashes(input, hash)

	// Assert
	if !isSameHash {
		t.Errorf("Expected true, but got false")
	}
}
