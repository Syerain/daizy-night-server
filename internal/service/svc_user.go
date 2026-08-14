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
	repoUser    abstract.InterfaceRepoUser
	repoToken   abstract.InterfaceRepoToken
	repoRegcode abstract.InterfaceRepoRegistercode
	crypto      abstract.InterfaceCrypto
}

func NewServiceUser(repoUser abstract.InterfaceRepoUser, repoToken abstract.InterfaceRepoToken, repoRegcode abstract.InterfaceRepoRegistercode, crypto abstract.InterfaceCrypto) *ServiceUser {
	return &ServiceUser{
		repoUser:    repoUser,
		repoToken:   repoToken,
		repoRegcode: repoRegcode,
		crypto:      crypto,
	}
}

func (s *ServiceUser) Register(b model.RegisterBody) (*model.User, error) {
	slog.Info("service processing register ...")

	// repo record regcode
	if err := s.repoRegcode.Record(b.Registercode); err != nil {
		return nil, err
	}

	// salt psw
	passwordHash, err := utils.HashCreate(b.Password)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	//
	switch b.Registerway {
	case consts.RegisterLegacy:
		payload, err := s.crypto.AnalyzeRegistercode(b.Registercode)
		if err != nil {
			return nil, err
		} //it can only be ErrRegistercode

		if !s.crypto.VerifyRegistercodePayload(*payload) {
			return nil, errs.BuildErrRegistercode(errs.RegistercodeUnusableOutdated, 400)
		}

		err = s.repoUser.CreateUser(&model.User{
			Username:     b.Username,
			Nickname:     b.Nickname,
			PasswordHash: passwordHash,
			Registercode: b.Registercode,
		})

		if err != nil {
			//slog.Error("failure during creating user;" + err.Error())
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return nil, errs.BuildErrValidation(errs.ValidationKeyDuplicatedValue, 400, string(consts.ExprIndetermined), string(consts.ExprIndetermined))
			}
			return nil, err
		}

	case consts.RegisterGithub:
		return nil, &errs.ErrSupport{
			Type: errs.FeatureUnsupported,
			Http: http.StatusInternalServerError,
		}
	default:
		return nil, &errs.ErrUnknown{
			Type: errs.Undefined,
			Http: http.StatusInternalServerError,
		}
	}

	user, err := s.repoUser.GetUserByUsername(b.Username)
	if err != nil {
		return nil, err
	}
	if err := s.repoRegcode.Used(b.Registercode, true, user.ID); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *ServiceUser) Login(b model.LoginBody) (success bool, accessToken string, refreshToken string, err error) {
	switch b.Loginway {
	case consts.LoginLegacy:
		{
			if b.Username != "" {
				user, err := s.repoUser.GetUserByUsername(b.Username)
				// db error during GetUser
				if err != nil {
					errapp, ok := errs.Easx[*errs.ErrDbRecord](err)
					if ok {
						if errapp.Type == errs.DbRecordUsernameNotFound {
							return false, "", "", errs.BuildErrUserLogin(errs.UserLoginParamsIncorrect, http.StatusBadRequest, string(consts.ExprIndetermined))
						}
					}
					return false, "", "", err
				}

				matched, err := utils.HashVerify(b.Password, user.PasswordHash)
				if err != nil {
					return false, "", "", err
				}
				if !matched {
					return false, "", "", errs.BuildErrUserLogin(errs.UserLoginParamsIncorrect, http.StatusBadRequest, user.Username)
				}

				payloadAccessToken := model.JwtAccessTokenPayload{
					Uid:      user.ID,
					Username: user.Username,
					Role:     user.Role,
				}
				payloadRefreshToken := model.JwtRefreshTokenPayload{
					Uid:      user.ID,
					Username: user.Username,
					Role:     user.Role,
				}

				accessToken, err := s.crypto.SignAccessToken(payloadAccessToken)
				if err != nil {
					slog.Error(err.Error())
					return false, "", "", err
				}

				refreshToken, err := s.crypto.SignRefreshToken(payloadRefreshToken)
				if err != nil {
					slog.Error(err.Error())
					return false, "", "", err
				}

				if err := s.repoToken.SaveRefreshToken(user.ID, refreshToken); err != nil {
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
	payload, err := s.crypto.VerifyRefreshToken(rawToken)
	if err != nil {
		return false, "", "", err
	}

	valid, err := s.repoToken.GetRefreshToken(payload.Uid, rawToken)
	if err != nil {
		return false, "", "", err
	}

	if !valid {
		return false, "", "", errs.BuildErrUserLogin(errs.UserLoginParamsIncorrect, http.StatusBadRequest, string(consts.ExprIndetermined))
	}

	if err = s.repoToken.RevokeUserTokens(payload.Uid); err != nil {
		return false, "", "", err
	}

	payloadAccess := model.JwtAccessTokenPayload{
		Uid:      payload.Uid,
		Username: payload.Username,
		Role:     payload.Role,
	}
	accessToken, err = s.crypto.SignAccessToken(payloadAccess)
	if err != nil {
		return false, "", "", err
	}

	payloadRefresh := model.JwtRefreshTokenPayload{
		Uid:      payload.Uid,
		Username: payload.Username,
		Role:     payload.Role,
	}
	refreshToken, err = s.crypto.SignRefreshToken(payloadRefresh)
	if err != nil {
		return false, "", "", err
	}

	err = s.repoToken.SaveRefreshToken(payload.Uid, refreshToken)
	if err != nil {
		return false, "", "", err
	}
	return true, accessToken, refreshToken, nil
}

func (s *ServiceUser) GetInfoMineByUid(uid uint) (*model.InfoMe, error) {
	user, err := s.repoUser.GetUserByUid(uid)
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
		GithubID:     user.GithubID,
		GithubLogin:  user.GithubLogin,
	}, nil
}

func (s *ServiceUser) GetUserByUsername(name string) (*model.User, error) {
	user, err := s.repoUser.GetUserByUsername(name)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *ServiceUser) GetUserByUid(uid uint) (*model.User, error) {
	user, err := s.repoUser.GetUserByUid(uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}
