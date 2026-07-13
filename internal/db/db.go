package db

import (
	"daizynight/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return err
	}
	DB = db
	return nil
}
