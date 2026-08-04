package handler

import (
	"log/slog"
	"net/http"

	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/atomreforge/daizy-night-server/internal/errs"
	"github.com/atomreforge/daizy-night-server/internal/model"
	"github.com/labstack/echo/v5"
)

func (h *HandlerComplex) HandleLogin(ctx *echo.Context) error {
	b, err := Bind[model.LoginBody](ctx)
	if err != nil {
		slog.Error("during login username::" + b.Username + ";" + err.Error())
		return Respond(ctx, http.StatusInternalServerError, string(consts.ExprHttpInternalServerError))
	}

	slog.Info("got login request; username: " + b.Username)

	if err := ValidateLoginParams(b); err != nil {
		if errapp, ok := errs.Easx[abstract.InterfaceAppError](err); ok {
			slog.Error("during login username::" + b.Username + ";" + errapp.Error())
			return RespondCustom(ctx, errapp)
		}
		slog.Error("during login username::" + b.Username + ";" + err.Error())
		return Respond(ctx, http.StatusInternalServerError, string(consts.ExprHttpInternalServerError))
	}

	ok, accessToken, refreshToken, err := h.ServiceUser.Login(*b)
	if !ok {
		if errapp, ok := errs.Easx[abstract.InterfaceAppError](err); ok {
			slog.Error("during login username::" + b.Username + ";" + errapp.Error())
			return RespondCustom(ctx, errapp)
		}
		slog.Error("during login username::" + b.Username + ";" + err.Error())

		return Respond(ctx, http.StatusInternalServerError, string(consts.ExprHttpInternalServerError))
	}

	slog.Info("successfully logged in user::" + b.Username)
	return RespondObj(ctx, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
