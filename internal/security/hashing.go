package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	random "math/rand"
)

type Hash struct {
	Hash string `json:"hash"`
	Salt string `json:"salt"`
}

func (h *Hash) String() string {
	return h.Hash + ":" + h.Salt
}

const ALPHABET string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func getRandomPepper() string {
	return string(ALPHABET[random.Int()%len(ALPHABET)])
}

func GenerateSalt(size int) []byte {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}
	return salt
}

func HashSHA256(input string) *Hash {
	salt := GenerateSalt(16)
	saltedInput := input + string(salt) + getRandomPepper()

	hash := sha256.New()
	hash.Write([]byte(saltedInput))
	hashed := hash.Sum(nil)

	finalHash := Hash{
		Hash: hex.EncodeToString(hashed),
		Salt: hex.EncodeToString(salt),
	}

	return &finalHash
}

func CompareHashes(input string, hash *Hash) bool {
	salt, _ := hex.DecodeString(hash.Salt)
	for i := 0; i < len(ALPHABET); i++ {
		saltedInput := input + string(salt) + string(ALPHABET[i])
		newHash := sha256.New()
		newHash.Write([]byte(saltedInput))
		newHashed := newHash.Sum(nil)
		if hex.EncodeToString(newHashed) == hash.Hash {
			return true
		}
	}
	return false
}
