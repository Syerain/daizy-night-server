package service

import (
	"daizynight/internal/db"
	"daizynight/internal/model"
	"daizynight/internal/utils"

	"daizynight/internal/constants"
	"log/slog"
)

func Register(b *model.RegisterBody) error {
	// handler is on duty of validation
	/* err := ValidateRegisterParams(b)
	if err != nil {
		return err
	} */

	slog.Info("service processing register ...")

	saltedPsw, err := utils.SaltMix(b.Password)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	switch b.Registerway.Value {
	case constants.Legacy:
		err = db.CreateUser(&model.User{
			Username:       b.Username,
			Nickname:       b.Nickname,
			SaltedPassword: saltedPsw,
			Registercode:   b.Registercode,
		})
		if err != nil {
			slog.Error(err.Error())
			return err
		}
	case constants.Github:
		return nil
	}

	return nil
}
