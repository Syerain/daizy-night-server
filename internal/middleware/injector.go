package middleware

import (
	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/atomreforge/daizy-night-server/internal/utils"

	"github.com/labstack/echo/v5"
)

func Inject() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			// trace id summon
			tid := ctx.Request().Header.Get("X-Request-ID")
			if tid == "" {
				tid = utils.NewTraceID()
			}

			// set addtional key: traceid, callChain, etc.
			ctx.Set(string(consts.ExprTraceid), tid)
			ctx.Set(string(consts.JsonExprCallChain), []string{})

			// json header
			ctx.Response().Header().Set("X-Request-ID", tid)

			return next(ctx)
		}
	}
}
