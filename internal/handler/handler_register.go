package handler

import (
	"daizynight/internal/model"
	"daizynight/internal/service"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
)

/* shit hill 's first show XD */

func (h *HandlerComplex) HandleRegister(ctx *echo.Context) error {
	b, err := Bind[model.RegisterBody](ctx)

	// failed to build RegisterBody
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "unknown register body eror"})
	}

	// failure during param validation
	err = ValidateRegisterParams(b)
	if err != nil {
		var errval *ErrRegisterValidation
		if errors.As(err, &errval) {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": errval.Message})
		}
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	// execute reg service
	if err := h.ServiceUser.Register(b); err != nil {
		var errReg *service.ErrRegister
		//
		if errors.As(err, &errReg) {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": errReg.Message})
		}
		// undefined err
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "internal error"})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "ok"})

}

func Bind[T any](ctx *echo.Context) (*T, error) {
	var b T
	if err := ctx.Bind(&b); err != nil {
		return nil, err
	}
	return &b, nil
}
