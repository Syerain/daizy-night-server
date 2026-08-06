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

var _ abstract.InterfaceServiceUser = (*ServiceUser)(nil)

type ServiceUser struct {
	RepoUser  abstract.InterfaceRepoUser
	RepoToken abstract.InterfaceRepoToken
	Crypto    abstract.InterfaceCrypto
}

func NewServiceUser(repouser abstract.InterfaceRepoUser, repotoken abstract.InterfaceRepoToken, crypto abstract.InterfaceCrypto) *ServiceUser {
	return &ServiceUser{
		RepoUser:  repouser,
		RepoToken: repotoken,
		Crypto:    crypto,
	}
}

func (s *ServiceUser) Register(b model.RegisterBody) error {
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
			return errs.BuildErrRegistercode(errs.RegistercodeUnusableOutdated, 400)
		}

		err = s.RepoUser.CreateUser(&model.User{
			Username:     b.Username,
			Nickname:     b.Nickname,
			PasswordHash: passwordHash,
			Registercode: b.Registercode,
		})

		if err != nil {
			//slog.Error("failure during creating user;" + err.Error())
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return errs.BuildErrValidation(errs.ValidationKeyDuplicatedValue, 400, string(consts.ExprIndetermined), string(consts.ExprIndetermined))
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
				user, err := s.RepoUser.GetUserByUsername(b.Username)
				// db error during GetUser
				if err != nil {
					return false, "", "", err
				}

				matched, err := utils.HashVerify(b.Password, user.PasswordHash)
				if err != nil {
					return false, "", "", err
				}
				if !matched {
					return false, "", "", errs.BuildErrUserLogin(errs.UserLoginParamsPasswordIncorrect, http.StatusBadRequest, user.Username)
				}

				payloadAccessToken := model.JwtAccessTokenPayload{
					Uid:      user.ID,
					Username: user.Username,
					Role:     user.Role,
				}
				payloadRefreshToken := model.JwtRefreshTokenPayload{
					Uid:      user.ID,
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

				if err := s.RepoToken.SaveRefreshToken(user.ID, refreshToken); err != nil {
					slog.Error(err.Error())
					return false, "", "", err
				}

				return true, accessToken, refreshToken, nil

			}
			return false, "", "", errs.BuildErrValidation(errs.ValidationKeyNull, http.StatusBadRequest, string(consts.ExprUsername), string(consts.ExprNull))
		}
	// under constructing ..
	case consts.LoginGithub:
		{
			return false, "", "", errs.BuildErrSupport(errs.FeatureUnsupported, http.StatusBadRequest)
		}
	}

	// wont arrive here. (hope so?)
	return false, "", "", errs.BuildErrSupport(errs.FeatureUnsupported, http.StatusBadRequest)
}

func (s *ServiceUser) RefreshAccessToken(rawToken string) (success bool, accessToken string, refreshToken string, err error) {
	payload, err := s.Crypto.VerifyRefreshToken(rawToken)
	if err != nil {
		return false, "", "", err
	}

	valid, err := s.RepoToken.GetRefreshToken(payload.Uid, rawToken)
	if err != nil || !valid {
		return false, "", "", err
	}

	err = s.RepoToken.RevokeUserTokens(payload.Uid)
	if err != nil {
		return false, "", "", err
	}

	payloadAccess := model.JwtAccessTokenPayload{
		Uid:      payload.Uid,
		Username: payload.Username,
	}
	accessToken, err = s.Crypto.SignAccessToken(payloadAccess)
	if err != nil {
		return false, "", "", err
	}

	payloadRefresh := model.JwtRefreshTokenPayload{
		Uid:      payload.Uid,
		Username: payload.Username,
	}
	refreshToken, err = s.Crypto.SignRefreshToken(payloadRefresh)
	if err != nil {
		return false, "", "", err
	}

	err = s.RepoToken.SaveRefreshToken(payload.Uid, refreshToken)
	if err != nil {
		return false, "", "", err
	}
	return true, accessToken, refreshToken, nil
}

func (s *ServiceUser) GetInfoMineByUid(uid uint) (*model.InfoMe, error) {
	user, err := s.RepoUser.GetUserByUid(uid)
	if err != nil {
		return nil, err
	}

	return &model.InfoMe{
		Uid:          user.ID,
		Username:     user.Username,
		Nickname:     user.Nickname,
		Email:        user.Email,
		RegisterTime: user.RegisterTime,
		Role:         user.Role,
		GitHubID:     user.GitHubID,
		GitHubLogin:  user.GitHubLogin,
	}, nil
}
