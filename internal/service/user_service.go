package service

import (
	"daizynight/internal/crypto"
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
	passwordHash, err := utils.HashCreate(b.Password)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	//
	switch b.Registerway {
	case constants.RegisterLegacy:
		err = db.CreateUser(&model.User{
			Username:     b.Username,
			Nickname:     b.Nickname,
			PasswordHash: passwordHash,
			Registercode: b.Registercode,
		})
		if err != nil {
			slog.Error(err.Error())
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return &ErrRegister{Message: "duplicated register params already exists in the database "}
			}
			return err
		}
	case constants.RegisterGithub:
		return nil
	default:
		panic("unreachable")
	}

	return nil
}

func Login(b model.LoginBody) (success bool, accessToken string, refreshToken string, err error) {
	switch b.Loginway {
	case constants.LoginLegacy:
		{
			if b.Username != "" {
				user, err := db.GetUserByUsername(b.Username)
				// db error during GetUser
				if err != nil {
					return false, "", "", err
				}
				matched, err := crypto.ValidateSaltedPassword(b.Password, user.PasswordHash)
				if err != nil {
					return false, "", "", err
				}
				if !matched {
					return false, "", "", nil
				}

				payloadAccessToken := crypto.JwtAccessTokenPayload{
					AtomID:   user.AtomID,
					Username: user.Username,
					Role:     user.Role,
				}
				payloadRefreshToken := crypto.JwtRefreshTokenPayload{
					AtomID:   user.AtomID,
					Username: user.Username,
				}

				accessToken, err := crypto.SignAccessToken(payloadAccessToken)
				if err != nil {
					slog.Error(err.Error())
					return false, "", "", err
				}

				refreshToken, err := crypto.SignRefreshToken(payloadRefreshToken)
				if err != nil {
					slog.Error(err.Error())
					return false, "", "", err
				}

				if err := db.SaveRefreshToken(user.AtomID, refreshToken); err != nil {
					slog.Error(err.Error())
					return false, "", "", err
				}

				return true, accessToken, refreshToken, nil

			}
		}
	// under constructing ..
	case constants.LoginGithub:
		{
			return false, "", "", nil
		}
	}

	// wont arrive here. (hope so?)
	return false, "", "", nil
}

func RefreshAccessToken(rawToken string) (success bool, accessToken string, refreshToken string, err error) {
	payload, err := crypto.VerifyRefreshToken(rawToken)
	if err != nil {
		return false, "", "", err
	}

	valid, err := db.GetRefreshToken(payload.AtomID, rawToken)
	if err != nil || !valid {
		return false, "", "", err
	}

	db.RevokeUserTokens(payload.AtomID)

	payloadAccess := crypto.JwtAccessTokenPayload{
		AtomID:   payload.AtomID,
		Username: payload.Username,
	}
	accessToken, err = crypto.SignAccessToken(payloadAccess)
	if err != nil {
		return false, "", "", err
	}

	payloadRefresh := crypto.JwtRefreshTokenPayload{
		AtomID:   payload.AtomID,
		Username: payload.Username,
	}
	refreshToken, err = crypto.SignRefreshToken(payloadRefresh)
	if err != nil {
		return false, "", "", err
	}

	db.SaveRefreshToken(payload.AtomID, refreshToken)
	return true, accessToken, refreshToken, nil
}
