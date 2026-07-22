package service

import (
	"daizynight/internal/db"
	"daizynight/internal/model"
	"daizynight/internal/utils"
	"errors"

	"daizynight/internal/constants"
	"log/slog"

	"gorm.io/gorm"
)

func Register(b *model.RegisterBody) error {
	slog.Info("service processing register ...")

	// salt psw
	saltedPsw, err := utils.SaltMix(b.Password)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	//
	switch b.Registerway {
	case constants.Legacy:
		err = db.CreateUser(&model.User{
			Username:       b.Username,
			Nickname:       b.Nickname,
			SaltedPassword: saltedPsw,
			Registercode:   b.Registercode,
		})
		if err != nil {
			slog.Error(err.Error())
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return &ErrRegister{Message: "duplicated register params already exists in the database "}
			}
			return err
		}
	case constants.Github:
		return nil
	}

	return nil
}
