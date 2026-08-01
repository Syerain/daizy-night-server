package service

import (
	abstract "daizynight/internal/abstract/interface"
	"daizynight/internal/constants"
	"daizynight/internal/crypto"
	"daizynight/internal/model"
	"daizynight/internal/utils"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type ServiceUser struct {
	Userrepo  abstract.InterfaceUserRepo
	Tokenrepo abstract.InterfaceTokenRepo
	Crypto    abstract.InterfaceCrypto
}

func NewServiceUser(repo abstract.InterfaceUserRepo, tokenrepo abstract.InterfaceTokenRepo, crypto abstract.InterfaceCrypto) *ServiceUser {
	return &ServiceUser{
		Userrepo:  repo,
		Tokenrepo: tokenrepo,
		Crypto:    crypto,
	}
}

func (s *ServiceUser) Register(b *model.RegisterBody) error {
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
		payload, err := s.Crypto.AnalyzeRegistercode(b.Registercode)
		if err != nil {
			var regErr *crypto.ErrRegister
			if errors.As(err, &regErr) {
				return &ErrRegister{Message: "invalid registercode format"}
			}
			return err // unreachable
		}
		if !s.Crypto.VerifyRegistercodePayload(*payload) {
			return &ErrRegister{Message: "unusable registercode"}
		}

		err = s.Userrepo.CreateUser(&model.User{
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

func (s *ServiceUser) Login(b model.LoginBody) (success bool, accessToken string, refreshToken string, err error) {
	switch b.Loginway {
	case constants.LoginLegacy:
		{
			if b.Username != "" {
				user, err := s.Userrepo.GetUserByUsername(b.Username)
				// db error during GetUser
				if err != nil {
					return false, "", "", err
				}
				matched, err := utils.HashVerify(b.Password, user.PasswordHash)
				if err != nil {
					return false, "", "", err
				}
				if !matched {
					return false, "", "", nil
				}

				payloadAccessToken := model.JwtAccessTokenPayload{
					AtomID:   user.AtomID,
					Username: user.Username,
					Role:     user.Role,
				}
				payloadRefreshToken := model.JwtRefreshTokenPayload{
					AtomID:   user.AtomID,
					Username: user.Username,
				}

				accessToken, err := s.Crypto.SignAccessToken(payloadAccessToken)
				if err != nil {
					slog.Error(err.Error())
					return false, "", "", err
				}

				refreshToken, err := s.Crypto.SignRefreshToken(payloadRefreshToken)
				if err != nil {
					slog.Error(err.Error())
					return false, "", "", err
				}

				if err := s.Tokenrepo.SaveRefreshToken(user.AtomID, refreshToken); err != nil {
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

func (s *ServiceUser) RefreshAccessToken(rawToken string) (success bool, accessToken string, refreshToken string, err error) {
	payload, err := s.Crypto.VerifyRefreshToken(rawToken)
	if err != nil {
		return false, "", "", err
	}

	valid, err := s.Tokenrepo.GetRefreshToken(payload.AtomID, rawToken)
	if err != nil || !valid {
		return false, "", "", err
	}

	err = s.Tokenrepo.RevokeUserTokens(payload.AtomID)
	if err != nil {
		return false, "", "", err
	}

	payloadAccess := model.JwtAccessTokenPayload{
		AtomID:   payload.AtomID,
		Username: payload.Username,
	}
	accessToken, err = s.Crypto.SignAccessToken(payloadAccess)
	if err != nil {
		return false, "", "", err
	}

	payloadRefresh := model.JwtRefreshTokenPayload{
		AtomID:   payload.AtomID,
		Username: payload.Username,
	}
	refreshToken, err = s.Crypto.SignRefreshToken(payloadRefresh)
	if err != nil {
		return false, "", "", err
	}

	err = s.Tokenrepo.SaveRefreshToken(payload.AtomID, refreshToken)
	if err != nil {
		return false, "", "", err
	}
	return true, accessToken, refreshToken, nil
}
