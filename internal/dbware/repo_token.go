package dbware

import (
	"time"

	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/model"
	"github.com/atomreforge/daizy-night-server/internal/utils"
	"gorm.io/gorm"
)

var _ abstract.InterfaceRepoToken = (*RepoToken)(nil)

type RepoToken struct {
	db *gorm.DB
}

func NewRepoToken(db *gorm.DB) *RepoToken {
	return &RepoToken{db: db}
}

func (r *RepoToken) SaveRefreshToken(uid uint, rawToken string) error {
	hash, err := utils.HashCreate(rawToken)
	if err != nil {
		return err
	}
	return r.db.Create(&model.RefreshToken{
		Uid:       uid,
		TokenHash: hash,
	}).Error
}

func (r *RepoToken) GetRefreshToken(uid uint, rawToken string) (bool, error) {
	var tokens []model.RefreshToken
	result := r.db.Where("id = ? AND revoked_at IS NULL", uid).Find(&tokens)
	if result.Error != nil {
		return false, result.Error
	}
	for _, t := range tokens {
		matched, err := utils.HashVerify(rawToken, t.TokenHash)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

func (r *RepoToken) RevokeUserTokens(uid uint) error {
	now := time.Now()
	return r.db.Model(&model.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", uid).
		Update("revoked_at", now).Error
}
