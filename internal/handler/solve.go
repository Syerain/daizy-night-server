package handler

import (
	"github.com/atomreforge/daizy-night-server/internal/model"
	"github.com/labstack/echo/v5"
)

func Bind[T any](ctx *echo.Context) (T, error) {
	var b T
	if err := ctx.Bind(&b); err != nil {
		return b, err
	}
	return b, nil
}

// requireSelf ensures the :username path segment denotes the authenticated
// user themself; any other username is forbidden (403).
func requireSelf(ctx *echo.Context, claims *model.JwtAccessTokenPayload) error {
	if ctx.Param("username") != claims.Username {
		return echo.ErrForbidden
	}
	return nil
}
