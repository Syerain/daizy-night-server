package router

import (
	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/handler"
	mid "github.com/atomreforge/daizy-night-server/internal/middleware"

	"github.com/labstack/echo/v5"
)

func New(
	h *handler.HandlerComplex,
	pCrypto abstract.InterfaceCrypto,
) *echo.Echo {
	e := echo.New()

	e.POST("/api/register", h.HandleRegister)
	e.POST("/api/login", h.HandleLogin)

	// under protection
	auth := e.Group("/api")
	auth.Use(mid.AuthenJWT(pCrypto))

	auth.GET("/me", h.HandleMe)

	return e
}
