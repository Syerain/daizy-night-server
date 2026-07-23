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
