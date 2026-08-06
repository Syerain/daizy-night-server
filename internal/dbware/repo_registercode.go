package dbware

import (
	"gorm.io/gorm"

	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/model"
)

var _ abstract.InterfaceRepoRegistercode = (*RepoRegistercode)(nil)

type RepoRegistercode struct {
	db *gorm.DB
}

func NewRepoRegistercode(db *gorm.DB) *RepoRegistercode {
	return &RepoRegistercode{db: db}
}

func (r *RepoRegistercode) Record(rawHex model.RegistercodeRawHex) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&model.RegistercodeRecord{RawHex: rawHex}).Error
	})
}

func (r *RepoRegistercode) Remove(rawHex model.RegistercodeRawHex) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Where("raw_hex = ?", string(rawHex)).
			Delete(&model.RegistercodeRecord{}).Error
	})
}
