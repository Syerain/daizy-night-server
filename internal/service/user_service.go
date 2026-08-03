package service

import (
	"errors"
	"log/slog"
	"net/http"

	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/consts"

	//"github.com/atomreforge/daizy-night-server/internal/crypto"
	"github.com/atomreforge/daizy-night-server/internal/errs"
	"github.com/atomreforge/daizy-night-server/internal/model"
	"github.com/atomreforge/daizy-night-server/internal/utils"

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
	case consts.RegisterLegacy:
		payload, err := s.Crypto.AnalyzeRegistercode(b.Registercode)
		if err != nil {
			return err
		} //it can only be ErrRegistercode

		if !s.Crypto.VerifyRegistercodePayload(*payload) {
			return &errs.ErrRegistercode{
				Type: errs.RegistercodeUnusableOutdated,
				Http: 400,
			}
		}

		err = s.Userrepo.CreateUser(&model.User{
			Username:     b.Username,
			Nickname:     b.Nickname,
			PasswordHash: passwordHash,
			Registercode: b.Registercode,
		})

		if err != nil {
			slog.Error("failure during creating user;" + err.Error())
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return &errs.ErrValidation{
					Type:  errs.ValidationKeyDuplicatedValue,
					Http:  400,
					Field: string(consts.ExprIndetermined),
					Value: string(consts.ExprIndetermined),
				}
			}
			return err
		}

	case consts.RegisterGithub:
		slog.Warn(string(consts.ExprUnsupportedFeature))
		return nil
	default:
		panic(string(consts.ExprUnreachableCase))
	}

	return nil
}

func (s *ServiceUser) Login(b model.LoginBody) (success bool, accessToken string, refreshToken string, err error) {
	switch b.Loginway {
	case consts.LoginLegacy:
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
					return false, "", "", &errs.ErrUserLogin{
						Type: errs.UserLoginParamsPasswordIncorrect,
						Http: http.StatusBadRequest,
					}
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
			return false, "", "", &errs.ErrValidation{
				Type:  errs.ValidationKeyNull,
				Http:  http.StatusBadRequest,
				Field: string(consts.ExprUsername),
				Value: string(consts.ExprNull),
			}
		}
	// under constructing ..
	case consts.LoginGithub:
		{
			return false, "", "", &errs.ErrSupport{
				Type: errs.FeatureUnsupported,
				Http: http.StatusBadRequest,
			}
		}
	}

	// wont arrive here. (hope so?)
	return false, "", "", &errs.ErrSupport{
		Type: errs.FeatureUnsupported,
		Http: http.StatusBadRequest,
	}
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
