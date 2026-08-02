package handler

import (
	"net/http"

	"github.com/atomreforge/daizy-night-server/internal/model"
	"github.com/atomreforge/daizy-night-server/internal/service"
	//"github.com/gin-gonic/gin"
)

func (h *HandlerComplex) HandleLogin(ctx *echo.Context) error {
	b, err := Bind[model.LoginBody](ctx)