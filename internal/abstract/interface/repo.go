package abstract

import (
	"github.com/atomreforge/daizy-night-server/internal/model"
)

type InterfaceUserRepo interface {
	CreateUser(b *model.User) error
	GetUserByUsername(name string) (*model.User, error)
	GetUserByAtomid(id uint) (*model.User, error)
}

type InterfaceTokenRepo interface {
	SaveRefreshToken(atomid uint, rawToken string) error
	GetRefreshToken(atomid uint, rawToken string) (bool, error)
	RevokeUserTokens(atomid uint) error
}
