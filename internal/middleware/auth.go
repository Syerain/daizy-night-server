package middleware

import (
	"crypto/ed25519"

	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/atomreforge/daizy-night-server/internal/crypto"
	"github.com/atomreforge/daizy-night-server/internal/model"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func AuthenJWT(p abstract.InterfaceCrypto) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContextKey: string(consts.ExprContextKeyJWT),
		KeyFunc:    crypto.NewJWTKeyFunc(p.GetJwtAccessTokenDeckey().(ed25519.PublicKey)),
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return &model.JwtAccessTokenPayload{}
		},
	})
}

/*//abandoned
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
}*/

func logPrintError(ctx *echo.Context, err error) string {
	return ("error::" + err.Error() +
		"method::" + ctx.Request().Method +
		"path::" + ctx.Request().URL.Path)
}
