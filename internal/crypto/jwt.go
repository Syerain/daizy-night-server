package crypto

import (
	"daizynight/internal/config"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessTokenExpireTime time.Duration
var refreshTokenExpireTime time.Duration

func InitJwt(cfg *config.Config) {
	accessTokenExpireTime = cfg.Http.JwtAccessTokenExpireTime
	refreshTokenExpireTime = cfg.Http.JwtRefreshTokenExpireTime
}

type AccessTokenPayload struct {
	// business
	Atomid   int    `json:"atomid"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isadmin"`

	// jwt attrs
	jwt.RegisteredClaims
}

type RefreshTokenPayload struct {
	Atomid   int    `json:"atomid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func SignAccessToken(payload AccessTokenPayload) (string, error) {
	payload.ExpiresAt = jwt.NewNumericDate(time.Now().Add(accessTokenExpireTime))
	payload.IssuedAt = jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, payload)
	return token.SignedString(jwtAccessTokenEnckey)
}

func SignRefreshToken(payload RefreshTokenPayload) (string, error) {
	payload.ExpiresAt = jwt.NewNumericDate(time.Now().Add(refreshTokenExpireTime))
	payload.IssuedAt = jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, payload)
	return token.SignedString(jwtAccessTokenEnckey)
}

func VerifyAccessToken(tokenStr string) (*AccessTokenPayload, error) {
	var payload AccessTokenPayload
	token, err := jwt.ParseWithClaims(tokenStr, &payload, func(t *jwt.Token) (interface{}, error) {
		return jwtAccessTokenDeckey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return &payload, nil
}

func VerifyRefreshToken(tokenStr string) (*RefreshTokenPayload, error) {
	var payload RefreshTokenPayload
	token, err := jwt.ParseWithClaims(tokenStr, &payload, func(t *jwt.Token) (interface{}, error) {
		return jwtRefreshTokenDeckey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return &payload, nil
}
