package crypto

import (
	"crypto/ed25519"
	"daizynight/internal/config"
	"encoding/hex"
	"time"
)

// byte[] here.
var registercodeEnckey ed25519.PrivateKey
var registercodeDeckey ed25519.PublicKey
var passwordEnckey ed25519.PrivateKey
var passwordDeckey ed25519.PublicKey
var accessTokenEnckey ed25519.PrivateKey
var accessTokenDeckey ed25519.PublicKey

// transfer hex into bytes; set global keys.
func Init(cfg *config.Config) error {
	var err error

	registercodeEnckey, err = hexToPrivKey(cfg.Security.RegistercodeEnckey)
	if err != nil {
		return err
	}
	registercodeDeckey, err = hexToPubKey(cfg.Security.RegistercodeDeckey)
	if err != nil {
		return err
	}

	passwordEnckey, err = hexToPrivKey(cfg.Security.PasswordEnckey)
	if err != nil {
		return err
	}
	passwordDeckey, err = hexToPubKey(cfg.Security.PasswordDeckey)
	if err != nil {
		return err
	}

	accessTokenEnckey, err = hexToPrivKey(cfg.Security.AccessTokenEnckey)
	if err != nil {
		return err
	}
	accessTokenDeckey, err = hexToPubKey(cfg.Security.AccessTokenDeckey)
	if err != nil {
		return err
	}

	return nil
}

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
func ValidateSignature(signature string) bool {
	if signature == "" {
		return false
	}
	dateStr := time.Now().Format("20060102") // havent use
	dateBytes := []byte(dateStr)
	sigBytes := []byte(signature)
	return ed25519.Verify(registercodeDeckey, dateBytes, sigBytes)
}
