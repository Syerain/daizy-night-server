package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/alexedwards/argon2id"
)

func HashCreate(raw string) (string, error) {
	return argon2id.CreateHash(raw, argon2id.DefaultParams)
}

func HashVerify(raw string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(raw, hash)
}

func SHA256HashHex(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
