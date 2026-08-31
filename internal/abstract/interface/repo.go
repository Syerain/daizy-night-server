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
	GetRefreshToken(uid uint, rawToken string) (*model.RefreshToken, error)
	// notice that this func will revoke all tokens of a user
	RevokeUserTokens(uid uint) error
	RotateRefreshToken(uid uint, usedLookupHash, newRawToken string) error
	GetUidByRefreshToken(rawToken string) (uint, error)
}

type InterfaceRepoRegistercode interface {
	//Record(registercodeRaw model.RegistercodeRawHex) error
	RecordNewRegistercode(b *model.RegistercodeRecord) error
	Remove(registercodeRaw model.RegistercodeRawHex) error
	Used(registercodRawHex model.RegistercodeRawHex, value bool, userID uint) error
	Updates(record model.RegistercodeRecord) error
	GetRecordByRegistercode(rawHex model.RegistercodeRawHex) (*model.RegistercodeRecord, error)
}
