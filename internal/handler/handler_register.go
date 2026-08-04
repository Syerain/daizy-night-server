package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/atomreforge/daizy-night-server/internal/consts"
	mid "github.com/atomreforge/daizy-night-server/internal/middleware"
	"github.com/atomreforge/daizy-night-server/internal/model"

	"github.com/labstack/echo/v5"
)

/* shit hill 's first show XD */

func (h *HandlerComplex) HandleRegister(ctx *echo.Context) error {
	slog.Info("got register request")
	b, err := Bind[model.RegisterBody](ctx)

	// failed to build RegisterBody
	if err != nil {
		return err
	}

	// failure during param validation
	if err := ValidateRegisterParams(b); err != nil {
		return err
		/*if errapp, ok := errs.Easx[abstract.InterfaceAppError](err); ok {
			return errs.BuildBadRequest(errapp.Error())
		}
		return errs.BuildInternal(string(consts.ExprHttpInternalServerError))*/
	}

	// execute reg service
	if err := h.ServiceUser.Register(b); err != nil {
		return err
		/*if errapp, ok := errs.Easx[abstract.InterfaceAppError](err); ok {
			return errs.BuildBadRequest(errapp.Error())
		}
		return errs.BuildInternal(string(consts.ExprHttpInternalServerError))*/
	}

	// process success
	//return ctx.JSON(http.StatusOK, map[string]string{"message": "ok"})
	u, err := h.ServiceUser.Userrepo.GetUserByUsername(b.Username)

	/*//corner case; wont appear.
	if err != nil {
		return err
	}*/
	slog.Info("succeeded register; user::" + b.Username + " atomid::" + fmt.Sprint(u.AtomID))
	return mid.Respond(ctx, http.StatusOK, string(consts.ExprHttpOk))
}

func Bind[T any](ctx *echo.Context) (*T, error) {
	var b T
	if err := ctx.Bind(&b); err != nil {
		return nil, err
	}
	return &b, nil
}
