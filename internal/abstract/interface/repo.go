package abstract

import (
	"github.com/atomreforge/daizy-night-server/internal/model"
)

type InterfaceRepoUser interface {
	CreateUser(b *model.User) error
	GetUserByUsername(name string) (*model.User, error)
	GetUserByUid(uid uint) (*model.User, error)
}

type InterfaceRepoToken interface {
	SaveRefreshToken(uid uint, rawToken string) error
	GetRefreshToken(uid uint, rawToken string) (bool, error)
	RevokeUserTokens(uid uint) error
}

type InterfaceRepoRegistercode interface {
	Record(registercodeRaw model.RegistercodeRawHex) error
	Remove(registercodeRaw model.RegistercodeRawHex) error
	Used(registercodRawHex model.RegistercodeRawHex, value bool) error
	Updates(record model.RegistercodeRecord) error
	GetRecordByRegistercode(rawHex model.RegistercodeRawHex) (*model.RegistercodeRecord, error)
}
