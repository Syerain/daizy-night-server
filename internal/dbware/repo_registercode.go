package dbware

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/atomreforge/daizy-night-server/internal/errs"
	"github.com/atomreforge/daizy-night-server/internal/model"
)

var _ abstract.InterfaceRepoRegistercode = (*RepoRegistercode)(nil)

type RepoRegistercode struct {
	pDB abstract.InterfaceProviderDB
}

func NewRepoRegistercode(pDB abstract.InterfaceProviderDB) *RepoRegistercode {
	return &RepoRegistercode{pDB: pDB}
}

func (r *RepoRegistercode) Record(rawHex model.RegistercodeRawHex) error {
	return r.pDB.DB().Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&model.RegistercodeRecord{RawHex: rawHex}).Error
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return &errs.ErrRegistercode{
					Type: errs.RegistercodeUnusableUsed,
					Http: http.StatusBadRequest,
				}
			}
			return err
		}
		return nil
	})
}

func (r *RepoRegistercode) Remove(rawHex model.RegistercodeRawHex) error {
	return r.pDB.DB().Transaction(func(tx *gorm.DB) error {
		return tx.Where("raw_hex = ?", string(rawHex)).
			Delete(&model.RegistercodeRecord{}).Error
	})
}

func (r *RepoRegistercode) Used(rawHex model.RegistercodeRawHex, value bool) error {
	rec, err := r.GetRecordByRegistercode(rawHex)
	if err != nil {
		return err
	}
	return r.pDB.DB().Transaction(func(tx *gorm.DB) error {
		return tx.Model(&rec).UpdateColumn("used", value).Error
	})
}

func (r *RepoRegistercode) Updates(record model.RegistercodeRecord) error {
	rec, err := r.GetRecordByRegistercode(record.RawHex)
	if err != nil {
		return err
	}
	record.ID = rec.ID

	return r.pDB.DB().Transaction(func(tx *gorm.DB) error {
		return tx.Model(&rec).Save(&record).Error
	})
}

func (r *RepoRegistercode) GetRecordByRegistercode(rawHex model.RegistercodeRawHex) (*model.RegistercodeRecord, error) {
	var rec model.RegistercodeRecord
	result := r.pDB.DB().Where("raw_hex = ?", string(rawHex)).First(&rec)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.BuildErrDbRecord(errs.DbRecordNotFound, http.StatusInternalServerError, string(consts.ExprRegistercodeRecord))
		}
		return nil, result.Error
	}
	return &rec, nil
}
