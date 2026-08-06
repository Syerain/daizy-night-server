package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	v1 "github.com/atomreforge/daizy-night-server/internal/api/v1"
	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/atomreforge/daizy-night-server/internal/errs"
	mid "github.com/atomreforge/daizy-night-server/internal/middleware"
	"gorm.io/gorm"

	"github.com/labstack/echo/v5"
)

func (h *HandlerComplex) HandleRegister(ctx *echo.Context) error {
	slog.Info("got register request")
	req, err := Bind[v1.RegisterRequest](ctx)

	// failed to build RegisterBody
	if err != nil {
		return err
	}

	// failure during param validation
	if err := ValidateRegisterParams(req.RegisterBody); err != nil {
		return err
	}

	// disable registercode; havent verify sig here;
	// bad regcode will be released in the error below.
	if err = h.ServiceCode.RecordRegistercode(req.Registercode); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.BuildErrRegistercode(errs.ValidationKeyDuplicatedValue, http.StatusBadRequest)
		}
		return err
	}

	// execute reg service
	if err := h.ServiceUser.Register(req.RegisterBody); err != nil {
		h.ServiceCode.RemoveRegistercode(req.Registercode)
		return err
	}

	// process success
	u, err := h.ServiceUser.RepoUser.GetUserByUsername(req.Username)

	//corner case; wont appear.
	if err != nil {
		if err0 := h.ServiceCode.RemoveRegistercode(req.Registercode); err0 != nil {
			// it doesnt matters; we dont throw it upward;
			slog.Error(string(consts.ExprFailedRegistercodeWithdraw))
		}
		return err
	}

	slog.Info("succeeded register; user::" + req.Username + " uid::" + fmt.Sprint(u.ID))
	return mid.Respond(ctx, http.StatusOK, string(consts.ExprHttpOk))
}

func Bind[T any](ctx *echo.Context) (*T, error) {
	var b T
	if err := ctx.Bind(&b); err != nil {
		return nil, err
	}
	return &b, nil
}
