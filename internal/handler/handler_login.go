package handler

import (
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
		return Respond(ctx, http.StatusInternalServerError, string(consts.ExprHttpInternalServerError))
	}

	if err := ValidateLoginParams(b); err != nil {
		if errapp, ok := errs.Easx[abstract.InterfaceAppError](err); ok {
			return RespondCustom(ctx, errapp)
		}
		return Respond(ctx, http.StatusInternalServerError, string(consts.ExprHttpInternalServerError))
	}

	ok, accessToken, refreshToken, err := h.ServiceUser.Login(*b)
	if !ok {
		if errapp, ok := errs.Easx[abstract.InterfaceAppError](err); ok {
			return RespondCustom(ctx, errapp)
		}
		return Respond(ctx, http.StatusInternalServerError, string(consts.ExprHttpInternalServerError))
	}

	return RespondObj(ctx, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
