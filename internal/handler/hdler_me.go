package handler

import (
	"log/slog"
	"net/http"

	v1 "github.com/atomreforge/daizy-night-server/internal/api/v1/user"
	"github.com/atomreforge/daizy-night-server/internal/consts"
	mid "github.com/atomreforge/daizy-night-server/internal/middleware"
	"github.com/atomreforge/daizy-night-server/internal/model"
	"github.com/atomreforge/daizy-night-server/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func (h *HandlerComplex) HandleMe(ctx *echo.Context) error {
	// record flow chain (monotonically accumulating)
	utils.AppendCallChain(ctx, string(consts.ModExprHandlerMe))
	utils.AppendCallChain(ctx, string(consts.ModExprServiceUser))

	token, err := echo.ContextGet[*jwt.Token](ctx, string(consts.ExprContextKeyJWT))
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(*model.JwtAccessTokenPayload)
	if !ok {
		slog.Error("failed to assert claims to JwtAccessTokenPayload")
		return echo.ErrUnauthorized
	}
	b, err := h.ServiceUser.GetInfoMineByUid(claims.Uid)
	if err != nil {
		return err
	}

	mid.RespondObj(ctx, http.StatusOK, v1.InfoMeResponse{InfoMe: *b})

	return nil

	/*token, err := echo.ContextGet[*jwt.Token](ctx, string(consts.ExprContextKeyJWT))
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(*model.JwtAccessTokenPayload)
	if !ok {
		slog.Error("failed to assert claims to JwtAccessTokenPayload")
		return echo.ErrUnauthorized
	}
	return mid.RespondObj(ctx, 200, map[string]any{
		string(consts.ExprUserID):   claims.Uid,
		string(consts.ExprUsername): claims.Username,
		string(consts.JsonExprRole): claims.Role,
	})*/
}
