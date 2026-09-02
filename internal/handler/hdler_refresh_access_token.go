package handler

import (
	"fmt"
	"net/http"

	v1 "github.com/atomreforge/daizy-night-server/internal/api/v1/user"
	"github.com/atomreforge/daizy-night-server/internal/consts"
	mid "github.com/atomreforge/daizy-night-server/internal/middleware"
	"github.com/atomreforge/daizy-night-server/internal/utils"
	"github.com/labstack/echo/v5"
)

func (h *HandlerComplex) HandleRefreshAccessToken(ctx *echo.Context) error {
	utils.AppendCallChain(ctx, string(consts.ModExprHandlerRefresh))
	utils.Layer(ctx).Info("got::" + string(consts.ExprReqRefresh))

	req, err := Bind[v1.RefrehAccessTokenRequest](ctx)
	if err != nil {
		return err
	}

	// a missing token is 400 validation, not 401 lookup-miss
	if err := ValidateRefreshParams(req.RefreshToken); err != nil {
		return err
	}

	// lookup twice; wait for optimization
	uid, err := h.ServiceUser.GetUidByRefreshToken(req.RefreshToken)
	if err != nil {
		return err
	}
	user, err := h.ServiceUser.GetUserByUid(uid)
	if err != nil {
		return err
	}

	ok, accessToken, refreshToken, err := h.ServiceUser.RefreshAccessToken(req.RefreshToken)
	if !ok || err != nil {
		return err
	}
	defer utils.Layer(ctx).Info(fmt.Sprintf("successfully refreshed access token for user::%s", user.Username))
	return mid.RespondObj(ctx, http.StatusOK, v1.RefreshAccessTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
