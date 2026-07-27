package utils

import (
	"github.com/alexedwards/argon2id"
)

func HashCreate(raw string) (string, error) {
	return argon2id.CreateHash(raw, argon2id.DefaultParams)
}

func HashVerify(raw string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(raw, hash)
}
