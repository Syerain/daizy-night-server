package utils

import (
	"github.com/alexedwards/argon2id"
)

func SaltMix(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return hash, &SaltGenError{Content: password, Message: err.Error()}
	}
	return hash, nil
}

func SaltVerify(password string, hash string) (bool, error) {
	isMatched, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return isMatched, &SaltVarifyError{Content: hash, Message: err.Error()}
	}
	if !isMatched {
		return false, nil
	}
	return true, nil
}
