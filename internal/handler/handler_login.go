package handler

import (
	"log/slog"
	"net/http"

	//abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	//"github.com/atomreforge/daizy-night-server/internal/consts"
	//"github.com/atomreforge/daizy-night-server/internal/errs"
	mid "github.com/atomreforge/daizy-night-server/internal/middleware"
	"github.com/atomreforge/daizy-night-server/internal/model"
	"github.com/labstack/echo/v5"
)

func (h *HandlerComplex) HandleLogin(ctx *echo.Context) error {
	b, err := Bind[model.LoginBody](ctx)
	if err != nil {
		return err
	}

	slog.Info("got login request; username:" + b.Username)

	if err := ValidateLoginParams(b); err != nil {
		return err
		/*if errapp, ok := errs.Easx[abstract.InterfaceAppError](err); ok {
			slog.Error("during login username:" + b.Username + ";" + errapp.Error())
		}
		slog.Error("during login username:" + b.Username + ";" + err.Error())
		return mid.Respond(ctx, http.StatusInternalServerError, string(consts.ExprHttpInternalServerError))*/
	}

	ok, accessToken, refreshToken, err := h.ServiceUser.Login(*b)
	if !ok {
		return err
		/*if errapp, ok := errs.Easx[abstract.InterfaceAppError](err); ok {
			slog.Error("during login username::" + b.Username + ";" + errapp.Error())
			return mid.RespondCustom(ctx, errapp)
		}
		slog.Error("during login username::" + b.Username + ";" + err.Error())

		return mid.Respond(ctx, http.StatusInternalServerError, string(consts.ExprHttpInternalServerError))*/
	}

	slog.Info("successfully logged in user::" + b.Username)
	return mid.RespondObj(ctx, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
