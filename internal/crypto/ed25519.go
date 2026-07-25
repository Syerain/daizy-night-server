package crypto

import (
	"crypto/ed25519"
	"daizynight/internal/config"
	"daizynight/internal/utils"
	"encoding/hex"
	"log/slog"
	"time"
)

// byte[] here.
var registercodeEnckey ed25519.PrivateKey
var registercodeDeckey ed25519.PublicKey
var passwordEnckey ed25519.PrivateKey
var passwordDeckey ed25519.PublicKey
var jwtAccessTokenEnckey ed25519.PrivateKey
var jwtAccessTokenDeckey ed25519.PublicKey
var jwtRefreshTokenEnckey ed25519.PrivateKey
var jwtRefreshTokenDeckey ed25519.PublicKey

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

	jwtAccessTokenEnckey, err = hexToPrivKey(cfg.Security.JwtAccessTokenEnckey)
	if err != nil {
		return err
	}
	jwtAccessTokenDeckey, err = hexToPubKey(cfg.Security.JwtAccessTokenDeckey)
	if err != nil {
		return err
	}

	jwtRefreshTokenEnckey, err = hexToPrivKey(cfg.Security.JwtRefreshTokenEnckey)
	if err != nil {
		return err
	}
	jwtRefreshTokenDeckey, err = hexToPubKey(cfg.Security.JwtRefreshTokenDeckey)
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
func ValidateRegistercode(code string) bool {
	if code == "" {
		return false
	}
	dateStr := time.Now().Format("20060102") // havent use
	dateBytes := []byte(dateStr)
	sigBytes := []byte(code)
	return ed25519.Verify(registercodeDeckey, dateBytes, sigBytes)
}

func ValidateSaltedPassword(psw string, saltedPsw string) (bool, error) {
	if psw == "" {
		return false, nil
	}
	//pswBytes := []byte(psw)
	//saltedPswBytes := []byte(saltedPsw)
	matched, err := utils.SaltVerify(psw, saltedPsw)
	if err != nil {
		slog.Error(err.Error())
		return false, err //argon2id internal error.
	}
	if !matched {
		return false, nil
	}
	return true, nil
}
