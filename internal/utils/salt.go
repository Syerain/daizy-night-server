package utils

import (
	"github.com/alexedwards/argon2id"
)

func SaltMix(raw string) (string, error) {
	hash, err := argon2id.CreateHash(raw, argon2id.DefaultParams)
	if err != nil {
		return hash, &SaltGenError{Content: raw, Message: err.Error()}
	}
	return hash, nil
}

func SaltVerify(raw string, hash string) (bool, error) {
	isMatched, err := argon2id.ComparePasswordAndHash(raw, hash)
	if err != nil {
		return isMatched, &SaltVarifyError{Content: hash, Message: err.Error()}
	}
	if !isMatched {
		return false, nil
	}
	return true, nil
}
