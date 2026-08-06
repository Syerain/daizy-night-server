package middleware

import (
	"net/http"

	"github.com/atomreforge/daizy-night-server/internal/config"
	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/labstack/echo/v5"
	echomw "github.com/labstack/echo/v5/middleware"
)

func RateLimit(cfg *config.Config) echo.MiddlewareFunc {
	rl := cfg.Http.RateLimit

	if !rl.Enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c *echo.Context) error {
				return next(c)
			}
		}
	}

	store := echomw.NewRateLimiterMemoryStoreWithConfig(
		echomw.RateLimiterMemoryStoreConfig{
			Rate:      rl.Rate,
			Burst:     rl.Burst,
			ExpiresIn: rl.ExpiresIn,
		},
	)

	return echomw.RateLimiterWithConfig(echomw.RateLimiterConfig{
		Skipper: echomw.DefaultSkipper,
		Store:   store,
		IdentifierExtractor: func(c *echo.Context) (string, error) {
			// unique client ip
			return c.RealIP(), nil
		},
		ErrorHandler: func(c *echo.Context, err error) error {
			return err
		},
		DenyHandler: func(c *echo.Context, identifier string, err error) error {
			return echo.NewHTTPError(http.StatusTooManyRequests, string(consts.HttpExprTooManyRequests))
		},
	})

}
