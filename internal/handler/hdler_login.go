package handler

import (
	"log/slog"
	"net/http"

	v1 "github.com/atomreforge/daizy-night-server/internal/api/v1"
	mid "github.com/atomreforge/daizy-night-server/internal/middleware"
	"github.com/labstack/echo/v5"
)

func (h *HandlerComplex) HandleLogin(ctx *echo.Context) error {
	req, err := Bind[v1.LoginRequest](ctx)
	if err != nil {
		return err
	}

	slog.Info("got login request; username:" + req.Username)

	if err := ValidateLoginParams(&req.LoginBody); err != nil {
		return err
	}

	ok, accessToken, refreshToken, err := h.ServiceUser.Login(req.LoginBody)
	if !ok {
		return err
	}

	slog.Info("successfully logged in user::" + req.Username)
	return mid.RespondObj(ctx, http.StatusOK, v1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
