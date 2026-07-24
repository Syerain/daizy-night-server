package crypto

import (
	"daizynight/internal/config"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

func InitJwt(cfg *config.Config) {

}

type AccessTokenPayload struct {
	// business
	UserID   uint   `json:"userid"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isadmin"`

	// jwt attrs
	jwt.RegisteredClaims
}

func SignAccessToken(payload AccessTokenPayload) (string, error) {
	payload.ExpiresAt = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	payload.IssuedAt = jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, payload)
	return token.SignedString(accessTokenEnckey)
}

func VerifyToken(tokenStr string) (*AccessTokenPayload, error) {
	var payload AccessTokenPayload
	token, err := jwt.ParseWithClaims(tokenStr, &payload, func(t *jwt.Token) (interface{}, error) {
		return accessTokenDeckey, nil
	})
	if err != nil || !token.Valid {
		return nil, jwt.ErrTokenInvalidId
	}
	return &payload, nil
}
