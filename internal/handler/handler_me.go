package handler

import (
	"log/slog"

	"github.com/atomreforge/daizy-night-server/internal/consts"
	mid "github.com/atomreforge/daizy-night-server/internal/middleware"
	"github.com/atomreforge/daizy-night-server/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func (h *HandlerComplex) HandleMe(ctx *echo.Context) error {
	token, err := echo.ContextGet[*jwt.Token](ctx, string(consts.ExprContextKeyJWT))
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(*model.JwtAccessTokenPayload)
	if !ok {
		slog.Error("failed to assert claims to JwtAccessTokenPayload")
		return echo.ErrUnauthorized
	}
	return mid.RespondObj(ctx, 200, map[string]any{
		string(consts.ExprAtomid):   claims.AtomID,
		string(consts.ExprUsername): claims.Username,
		string(consts.JsonExprRole): claims.Role,
	})
}
