package dbware

import (
	"daizynight/internal/model"
	"daizynight/internal/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

// gorm had packaged db features but more business functions is needed.
// dbware.go provides them.
// we dont wrap all the functions bcz that is an abstract layer over another.

func (p *ProviderDB) CreateUser(b *model.User) error {
	b.RegisterTime = time.Now()
	result := p.db.Create(b)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return result.Error
		}
		return result.Error
	}
	return nil
}

func (p *ProviderDB) GetUserByUsername(name string) (*model.User, error) {
	var user model.User
	result := p.db.Where(&model.User{Username: name}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
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
