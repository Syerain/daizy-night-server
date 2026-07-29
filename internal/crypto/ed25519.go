package crypto

import (
	"crypto/ed25519"
	"daizynight/internal/utils"
	"encoding/hex"
	"log/slog"
	"time"
)

// byte[] here.
var registercodeEnckey ed25519.PrivateKey
var registercodeDeckey ed25519.PublicKey
var jwtAccessTokenEnckey ed25519.PrivateKey
var jwtAccessTokenDeckey ed25519.PublicKey
var jwtRefreshTokenEnckey ed25519.PrivateKey
var jwtRefreshTokenDeckey ed25519.PublicKey

func hexToPrivKey(s string) (ed25519.PrivateKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return ed25519.PrivateKey(b), nil
}

func hexToPubKey(s string) (ed25519.PublicKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return ed25519.PublicKey(b), nil
}

// 1:correct key; 2.correct date
func ValidateRegistercode(code string) bool {
	if code == "" {
		return false
	}
	dateStr := time.Now().Format("20060102") // havent use
	dateBytes := []byte(dateStr)
	sigBytes := []byte(code)
	return ed25519.Verify(registercodeDeckey, dateBytes, sigBytes)
}

type Password string

func (p *Password) Validate(password, hash string) error {
	// TODO
	return nil
}

func ValidateSaltedPassword(psw string, saltedPsw string) (bool, error) {
	//pswBytes := []byte(psw)
	//saltedPswBytes := []byte(saltedPsw)
	matched, err := utils.HashVerify(psw, saltedPsw)
	if err != nil {
		slog.Error(err.Error())
		return false, err //argon2id internal error.
	}
	return matched, nil
}
