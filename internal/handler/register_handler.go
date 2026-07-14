package handler

import (
	"daizynight/internal/model"
	"daizynight/internal/service"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
)

/* shit hill 's first show XD */

func HandleRegister(c *echo.Context) error {
	b, err := buildRegisterBody(c)

	// failed to build RegisterBody
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}

	// execute reg service
	if err := service.Register(b); err != nil {
		var errReg *service.RegisterError
		if errors.As(err, &errReg) {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": errReg.Message})
		}
		// undefined err
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "ok"})

}

func buildRegisterBody(c *echo.Context) (*model.RegisterBody, error) {
	var b model.RegisterBody
	if err := c.Bind(&b); err != nil {
		return nil, err
	}
	return &b, nil
}
