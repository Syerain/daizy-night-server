package dbware

import (
	"errors"
	"net/http"
	"time"

	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/atomreforge/daizy-night-server/internal/errs"
	"github.com/atomreforge/daizy-night-server/internal/model"
	"github.com/atomreforge/daizy-night-server/internal/utils"

	"gorm.io/gorm"
)

var _ abstract.InterfaceRepoUser = (*RepoUser)(nil)

type RepoUser struct {
	pDB abstract.InterfaceProviderDB
}

func NewRepoUser(pDB abstract.InterfaceProviderDB) *RepoUser {
	return &RepoUser{pDB: pDB}
}

/*func (r *RepoUser) CreateUser(b *model.User) error {
	b.RegisterTime = time.Now()

	return r.pDB.DB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(b).Error
	})
}*/

func (r *RepoUser) CreateUser(b *model.User) error {
	return r.pDB.DB().Transaction(func(tx *gorm.DB) error {
		uid, err := utils.GenUid()
		if err != nil {
			return err
		}
		b.UserID = uid
		b.RegisterTime = time.Now()

		/*latestID, err := r.GetLatestUid()
		if err != nil {
			return errs.BuildErrDbRecord(errs.Unknown, http.StatusInternalServerError, string(consts.ExprUserID))
		}
		b.UserID = latestID + 1*/

		if err := tx.Create(b).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return errs.BuildErrValidation(errs.ValidationKeyDuplicatedValue, http.StatusBadRequest, string(consts.ExprUsername), b.Username)
			}
			return err
		}

		err = tx.Create(&model.RegistercodeRecord{
			RawHex: b.Registercode,
			Used:   true,
			UserID: &b.UserID,
		}).Error
		if err == nil {
			return nil
		}
		if !errors.Is(err, gorm.ErrDuplicatedKey) {
			return err
		}

		res := tx.Model(&model.RegistercodeRecord{}).
			Where("raw_hex = ? AND used = ?", string(b.Registercode), false).
			Updates(map[string]any{"used": true, "user_id": b.UserID})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errs.BuildErrRegistercode(errs.RegistercodeUnusableUsed, http.StatusBadRequest)
		}
		return nil
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
	result := r.pDB.DB().Where("user_id = ?", uid).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.BuildErrDbRecord(errs.DbRecordNotFound, http.StatusBadRequest, string(consts.ExprUser))
		}
		return nil, result.Error
	}
	return &user, nil
}
