package handler

import (
	"net/http"

	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/atomreforge/daizy-night-server/internal/errs"
	"github.com/atomreforge/daizy-night-server/internal/model"

	"github.com/labstack/echo/v5"
)

/* shit hill 's first show XD */

func (h *HandlerComplex) HandleRegister(ctx *echo.Context) error {
	b, err := Bind[model.RegisterBody](ctx)

	// failed to build RegisterBody
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": string(consts.ExprHttpInternalServerError)})
	}

	// failure during param validation
	if err := ValidateRegisterParams(b); err != nil {
		if errapp, ok := errs.Easx[abstract.InterfaceAppError](err); ok {
			return ctx.JSON(errapp.HttpAbort(), map[string]string{"message": errapp.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": string(consts.ExprHttpInternalServerError)})
	}

	// execute reg service
	if err := h.ServiceUser.Register(b); err != nil {
		if errapp, ok := errs.Easx[abstract.InterfaceAppError](err); ok {
			return ctx.JSON(errapp.HttpAbort(), map[string]string{"message": errapp.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": string(consts.ExprHttpInternalServerError)})
	}

	// process success
	return ctx.JSON(http.StatusOK, map[string]string{"message": "ok"})

}

func Bind[T any](ctx *echo.Context) (*T, error) {
	var b T
	if err := ctx.Bind(&b); err != nil {
		return nil, err
	}
	return &b, nil
}
