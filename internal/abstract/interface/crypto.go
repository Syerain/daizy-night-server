package abstract

import (
	"daizynight/internal/model"
)

type InterfaceCrypto interface {
	SignRegistercode(payload model.RegistercodePayload) (string, error)
	AnalyzeRegistercode(codeStr string) (*model.RegistercodePayload, error)
	VerifyRegistercodePayload(payload model.RegistercodePayload) bool
	SignAccessToken(payload model.JwtAccessTokenPayload) (string, error)
	SignRefreshToken(payload model.JwtRefreshTokenPayload) (string, error)
	VerifyAccessToken(tokenStr string) (*model.JwtAccessTokenPayload, error)
	VerifyRefreshToken(tokenStr string) (*model.JwtRefreshTokenPayload, error)
}
