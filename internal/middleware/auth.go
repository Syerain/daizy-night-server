package middleware

import (
	"log/slog"
	"net/http"

	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/atomreforge/daizy-night-server/internal/errs"
	"github.com/labstack/echo/v5"
)

func ErrorRecovery() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("panic recovered",
						slog.Any("panic", r),
						slog.String("method", ctx.Request().Method),
						slog.String("path", ctx.Request().URL.Path),
					)
					Respond(ctx, http.StatusInternalServerError, string(consts.ExprHttpInternalServerError))
				}
			}()

			err := next(ctx)

			if err == nil {
				return nil
			}

			if res, ok := ctx.Response().(*echo.Response); ok && res.Committed {
				return err
			}

			if errapp, ok := errs.Easx[abstract.InterfaceAppError](err); ok {
				slog.Error("app error;" + logPrintError(ctx, errapp))
				return RespondCustom(ctx, errapp)
			}

			slog.Error("unhandled error" + logPrintError(ctx, err))
			return Respond(ctx, http.StatusInternalServerError, string(consts.ExprHttpInternalServerError))
		}
	}
}

func logPrintError(ctx *echo.Context, err error) string {
	return ("error::" + err.Error() +
		"method::" + ctx.Request().Method +
		"path::" + ctx.Request().URL.Path)
}
