package dbware

import (
	"time"

	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/config"
	"github.com/atomreforge/daizy-night-server/internal/model"
	"github.com/atomreforge/daizy-night-server/internal/utils"
)

var _ abstract.InterfaceRepoToken = (*RepoToken)(nil)

type RepoToken struct {
	pDB abstract.InterfaceProviderDB
	cfg *config.Config
}

func NewRepoToken(pDB abstract.InterfaceProviderDB, cfg *config.Config) *RepoToken {
	return &RepoToken{pDB: pDB, cfg: cfg}
}

func (r *RepoToken) SaveRefreshToken(uid uint, rawToken string) error {
	// clear expired tokens
	if err := r.PruneRevokedTokens(uid); err != nil {
		return err
	}

	hash, err := utils.HashCreate(rawToken)
	if err != nil {
		return err
	}
	return r.pDB.DB().Create(&model.RefreshToken{
		Uid:       uid,
		TokenHash: hash,
	}).Error
}

func (r *RepoToken) GetRefreshToken(uid uint, rawToken string) (bool, error) {
	var tokens []model.RefreshToken
	result := r.pDB.DB().Where("id = ? AND revoked_at IS NULL", uid).Find(&tokens)
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
	return r.pDB.DB().Model(&model.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", uid).
		Update("revoked_at", now).Error
}

func (r *RepoToken) PruneRevokedTokens(uid uint) error {
	cutoff := time.Now().Add(-r.cfg.Security.JwtRevokedTokensRetainTime)
	return r.pDB.DB().Where("uid = ? AND revoked_at IS NOT NULL AND revoked_at < ?", uid, cutoff).
		Delete(&model.RefreshToken{}).Error
}
