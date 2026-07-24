package db

import (
	"daizynight/internal/model"
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
	result := DB.Where(&model.User{Atomid: id}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
