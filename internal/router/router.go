package router

import (
	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/config"
	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/atomreforge/daizy-night-server/internal/handler"
	mid "github.com/atomreforge/daizy-night-server/internal/middleware"

	"github.com/labstack/echo/v5"
)

func New(
	h *handler.HandlerComplex,
	pCrypto abstract.InterfaceCrypto,
	cfg *config.Config,
) *echo.Echo {
	// init
	e := echo.New()

	e.Use(mid.Inject())       // injector
	e.Use(mid.RateLimit(cfg)) // rate limiter

	// public routes
	e.POST("/api/v1/register", h.HandleRegister)
	e.POST("/api/v1/login", h.HandleLogin)
	e.GET("/user/refresh-access-token", h.HandleRefreshAccessToken)

	// endpoints that requires authen only
	ptAuthOnly := e.Group("/api/v1")
	ptAuthOnly.Use(mid.AuthenJWT(pCrypto))
	ptAuthOnly.GET("/user/me", h.HandleMe)

	// admin only endpoints
	ptAdmin := e.Group("/api/v1/admin")
	ptAdmin.Use(mid.AuthenJWT(pCrypto))
	ptAdmin.Use(mid.RoleControl(consts.Admin))
	ptAdmin.POST("/sudo", h.HandleAdminSudo)

	return e
}
