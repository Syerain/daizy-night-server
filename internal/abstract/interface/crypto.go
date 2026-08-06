package abstract

import (
	"github.com/atomreforge/daizy-night-server/internal/model"
)

type InterfaceCrypto interface {
	SignRegistercode(payload model.RegistercodePayload) (model.RegistercodeRawHex, error)
	AnalyzeRegistercode(codeStr model.RegistercodeRawHex) (*model.RegistercodePayload, error)
	VerifyRegistercodePayload(payload model.RegistercodePayload) bool

	SignAccessToken(payload model.JwtAccessTokenPayload) (string, error)
	SignRefreshToken(payload model.JwtRefreshTokenPayload) (string, error)
	VerifyAccessToken(tokenStr string) (*model.JwtAccessTokenPayload, error)

	VerifyRefreshToken(tokenStr string) (*model.JwtRefreshTokenPayload, error)

	GetJwtAccessTokenDeckey() any
	GetJwtRefreshTokenDeckey() any
}
