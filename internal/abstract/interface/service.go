package abstract

import "github.com/atomreforge/daizy-night-server/internal/model"

type InterfaceServiceUser interface {
	Register(b model.RegisterBody) error
	Login(b model.LoginBody) (success bool, accessToken string, refreshToken string, err error)
	RefreshAccessToken(rawToken string) (success bool, accessToken string, refreshToken string, err error)
}

type InterfaceServiceCode interface {
	RecordRegistercode(registercodeRaw model.RegistercodeRawHex) error
	RemoveRegistercode(registercodeRaw model.RegistercodeRawHex) error
}

type InterfaceServiceAdmin interface {
	Sudo()
}
