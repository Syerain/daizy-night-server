package dbware

import (
	"errors"
	"net/http"
	"time"

	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/errs"
	"github.com/atomreforge/daizy-night-server/internal/model"
	"github.com/atomreforge/daizy-night-server/internal/utils"

	"gorm.io/gorm"
)

var (
	_ abstract.InterfaceUserRepo  = (*ProviderDB)(nil)
	_ abstract.InterfaceTokenRepo = (*ProviderDB)(nil)
)

// gorm had packaged db features but more business functions is needed.
// dbware.go provides them.
// we dont wrap all the functions bcz that is an abstract layer over another.

func (p *ProviderDB) CreateUser(b *model.User) error {
	b.RegisterTime = time.Now()

	return p.db.Transaction(func(tx *gorm.DB) error {
		maxid, err := p.getLatestAtomid()
		if err != nil {
			return err
		}
		b.AtomID = maxid + 1
		return tx.Create(b).Error
	})
}

/*// abandoned
func (p *ProviderDB) CreateUser(b *model.User) error {
	b.RegisterTime = time.Now()
	id, err := p.GetLatestAtomid()
	if err != nil {
		return err
	}
	b.AtomID = id + 1
	result := p.db.Create(b)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return result.Error
		}
		return result.Error
	}
	return nil
}*/

func (p *ProviderDB) GetUserByUsername(name string) (*model.User, error) {
	var user model.User
	result := p.db.Where(&model.User{Username: name}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.BuildErrDbRecord(errs.DbRecordUsernameNotFound, http.StatusBadRequest)
		}
		return nil, result.Error
	}
	return &user, nil
}

func (p *ProviderDB) GetUserByAtomid(atomid uint) (*model.User, error) {
	var user model.User
	result := p.db.Where(&model.User{AtomID: atomid}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (p *ProviderDB) SaveRefreshToken(atomid uint, rawToken string) error {
	hash, err := utils.HashCreate(rawToken)
	if err != nil {
		return err
	}
	return p.db.Create(&model.RefreshToken{
		AtomID:    atomid,
		TokenHash: hash,
	}).Error
}

func (p *ProviderDB) GetRefreshToken(atomid uint, rawToken string) (bool, error) {
	var tokens []model.RefreshToken
	result := p.db.Where("atomid = ? AND revoked_at IS NULL", atomid).Find(&tokens)
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

func (p *ProviderDB) RevokeUserTokens(atomid uint) error {
	now := time.Now()
	return p.db.Model(&model.RefreshToken{}).
		Where("atomid = ? AND revoked_at IS NULL", atomid).
		Update("revoked_at", now).Error
}

func (p *ProviderDB) getLatestAtomid() (uint, error) {
	var maxid uint
	err := p.db.Table("users").
		Select("COALESCE(MAX(atom_id), 0)").
		Scan(&maxid).Error
	return maxid, err
}

/*// abandoned due to async conflicts
func (p *ProviderDB) GetLatestAtomid() (uint, error) {
	var result struct {
		AtomID uint `gorm:"column:atomid"`
	}

	err := p.db.Table("users").
		Select("atomid").
		Order("atomid DESC").
		Limit(1).
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, &errs.ErrDbRecord{
			Type: errs.DbRecordNotFound,
			Http: http.StatusInternalServerError,
		}
	}
	return result.AtomID, nil
}*/
