package crypto

import (
	"crypto/ed25519"
	"errors"

	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/config"
	"github.com/atomreforge/daizy-night-server/internal/model"
	"github.com/atomreforge/daizy-night-server/internal/utils"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

var _ abstract.InterfaceCrypto = (*ProviderCrypto)(nil)

type ProviderCrypto struct {
	JwtAccessTokenEnckey   ed25519.PrivateKey
	JwtAccessTokenDeckey   ed25519.PublicKey
	JwtRefreshTokenEnckey  ed25519.PrivateKey
	JwtRefreshTokenDeckey  ed25519.PublicKey
	RegistercodeEnckey     ed25519.PrivateKey
	RegistercodeDeckey     ed25519.PublicKey
	AccessTokenExpireTime  time.Duration
	RefreshTokenExpireTime time.Duration
}

func NewProviderCrypto(cfg *config.Config) (*ProviderCrypto, error) {
	var errs []error
	jate, err := utils.HexToPrivKey(cfg.Security.JwtAccessTokenEnckey)
	errs = append(errs, err)
	jatd, err := utils.HexToPubKey(cfg.Security.JwtAccessTokenDeckey)
	errs = append(errs, err)
	jrte, err := utils.HexToPrivKey(cfg.Security.JwtRefreshTokenEnckey)
	errs = append(errs, err)
	jrtd, err := utils.HexToPubKey(cfg.Security.JwtRefreshTokenDeckey)
	errs = append(errs, err)
	re, err := utils.HexToPrivKey(cfg.Security.RegistercodeEnckey)
	errs = append(errs, err)
	rd, err := utils.HexToPubKey(cfg.Security.RegistercodeDeckey)
	errs = append(errs, err)
	atet := cfg.Security.JwtAccessTokenExpireTime
	rtet := cfg.Security.JwtRefreshTokenExpireTime

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return &ProviderCrypto{
		JwtAccessTokenEnckey:   jate,
		JwtAccessTokenDeckey:   jatd,
		JwtRefreshTokenEnckey:  jrte,
		JwtRefreshTokenDeckey:  jrtd,
		RegistercodeEnckey:     re,
		RegistercodeDeckey:     rd,
		AccessTokenExpireTime:  atet,
		RefreshTokenExpireTime: rtet,
	}, nil
}

func (p *ProviderCrypto) SignAccessToken(payload model.JwtAccessTokenPayload) (string, error) {
	payload.ExpiresAt = jwt.NewNumericDate(time.Now().Add(p.AccessTokenExpireTime))
	payload.IssuedAt = jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, payload)
	return token.SignedString(p.JwtAccessTokenEnckey)
}

func (p *ProviderCrypto) SignRefreshToken(payload model.JwtRefreshTokenPayload) (string, error) {
	payload.ExpiresAt = jwt.NewNumericDate(time.Now().Add(p.RefreshTokenExpireTime))
	payload.IssuedAt = jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, payload)
	return token.SignedString(p.JwtRefreshTokenEnckey)
}

func (p *ProviderCrypto) VerifyAccessToken(tokenStr string) (*model.JwtAccessTokenPayload, error) {
	var payload model.JwtAccessTokenPayload
	token, err := jwt.ParseWithClaims(tokenStr, &payload, func(t *jwt.Token) (interface{}, error) {
		return p.JwtAccessTokenDeckey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return &payload, nil
}

func (p *ProviderCrypto) VerifyRefreshToken(tokenStr string) (*model.JwtRefreshTokenPayload, error) {
	var payload model.JwtRefreshTokenPayload
	token, err := jwt.ParseWithClaims(tokenStr, &payload, func(t *jwt.Token) (interface{}, error) {
		return p.JwtRefreshTokenDeckey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return &payload, nil
}
