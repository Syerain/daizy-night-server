package db

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

func CreateUser(b *model.User) error {
	b.RegisterTime = time.Now()
	result := DB.Create(b)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return result.Error
		}
		return result.Error
	}
	return nil
}

func GetUserByUsername(name string) (*model.User, error) {
	var user model.User
	result := DB.Where(&model.User{Username: name}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByAtomid(id int) (*model.User, error) {
	var user model.User
	result := DB.Where(&model.User{AtomID: id}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func SaveRefreshToken(atomid int, rawToken string) error {
	hash, err := utils.HashCreate(rawToken)
	if err != nil {
		return err
	}
	return DB.Create(&model.RefreshToken{
		AtomID:    atomid,
		TokenHash: hash,
	}).Error
}

func GetRefreshToken(atomid int, rawToken string) (bool, error) {
	var tokens []model.RefreshToken
	result := DB.Where("atomid = ? AND revoked_at IS NULL", atomid).Find(&tokens)
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

func RevokeUserTokens(atomid int) error {
	now := time.Now()
	return DB.Model(&model.RefreshToken{}).
		Where("atomid = ? AND revoked_at IS NULL", atomid).
		Update("revoked_at", now).Error
}
