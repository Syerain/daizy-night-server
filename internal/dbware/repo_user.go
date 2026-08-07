package dbware

import (
	"errors"
	"net/http"
	"time"

	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/atomreforge/daizy-night-server/internal/errs"
	"github.com/atomreforge/daizy-night-server/internal/model"

	"gorm.io/gorm"
)

var _ abstract.InterfaceRepoUser = (*RepoUser)(nil)

type RepoUser struct {
	pDB abstract.InterfaceProviderDB
}

func NewRepoUser(pDB abstract.InterfaceProviderDB) *RepoUser {
	return &RepoUser{pDB: pDB}
}

func (r *RepoUser) CreateUser(b *model.User) error {
	b.RegisterTime = time.Now()

	return r.pDB.DB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(b).Error
	})
}

func (r *RepoUser) GetUserByUsername(name string) (*model.User, error) {
	var user model.User
	result := r.pDB.DB().Where(&model.User{Username: name}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.BuildErrDbRecord(errs.DbRecordUsernameNotFound, http.StatusBadRequest, string(consts.ExprUser))
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *RepoUser) GetUserByUid(uid uint) (*model.User, error) {
	var user model.User
	result := r.pDB.DB().Where("id = ?", uid).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.BuildErrDbRecord(errs.DbRecordNotFound, http.StatusBadRequest, string(consts.ExprUser))
		}
		return nil, result.Error
	}
	return &user, nil
}
