package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/config"
	"github.com/atomreforge/daizy-night-server/internal/consts"
	"github.com/atomreforge/daizy-night-server/internal/crypto"
	"github.com/atomreforge/daizy-night-server/internal/dbware"
	"github.com/atomreforge/daizy-night-server/internal/errs"
	"github.com/atomreforge/daizy-night-server/internal/handler"
	"github.com/atomreforge/daizy-night-server/internal/router"
	"github.com/atomreforge/daizy-night-server/internal/service"
	"github.com/atomreforge/daizy-night-server/internal/utils"

	mid "github.com/atomreforge/daizy-night-server/internal/middleware"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type App struct {
	echo *echo.Echo
	cfg  *config.Config
	db   *dbware.ProviderDB
}

func New(cfg *config.Config, dbProvider *dbware.ProviderDB) *App {
	// logger
	utils.InitModuleLogger(cfg.Main.IsDebugMode, "main")

	return &App{
		echo: echo.New(),
		cfg:  cfg,
		db:   dbProvider,
	}
}

func main() {
	// init
	slog.Info("Server starting...")
	slog.Info("Reading config...")
	cfg := config.MustLoadConfig()
	utils.InitModuleLogger(cfg.Main.IsDebugMode, "main")
	slog.Info("Debug Mode:", slog.Bool("enabled", cfg.Main.IsDebugMode))
	slog.Info("Configured listening at", "address", cfg.Http.Address, "port", cfg.Http.Port)

	// database provider
	db, err := dbware.NewDBProvider(struct {
		IsDebugMode bool
		DSN         string
	}{
		IsDebugMode: cfg.Main.IsDebugMode,
		DSN:         cfg.Database.DSN,
	})
	if err != nil {
		slog.Error("FATAL: couldnt init db !")
		os.Exit(1)
	}
	defer db.Close()

	// crypto provider
	cryptoj, err := crypto.NewProviderCrypto(cfg)
	if err != nil {
		slog.Error("FATAL: couldnt init crypto !")
		os.Exit(1)
	}

	// service provider
	svcUser := service.NewServiceUser(db, db, cryptoj)

	// handler
	handler := &handler.HandlerComplex{
		ServiceUser: svcUser,
	}

	// Echo
	e := router.New(handler)
	e.Logger = utils.GetLogger()
	e.HTTPErrorHandler = func(ctx *echo.Context, err error) {
		if resp, uErr := echo.UnwrapResponse(ctx.Response()); uErr == nil {
			if resp.Committed {
				return
			}
		}
		// set default to ISE
		code := http.StatusInternalServerError
		var sc echo.HTTPStatusCoder
		if errors.As(err, &sc) {
			if tmp := sc.StatusCode(); tmp != 0 {
				code = tmp //http status code; not biz code
				//slog.Error(err.Error())
				errapp, ok := errs.Easx[abstract.InterfaceAppError](err)
				if ok {
					mid.RespondCustom(ctx, errapp)
				} else {
					mid.Respond(ctx, code, string(consts.ExprHttpInternalServerError))
				}
			} else {
				mid.Respond(ctx, http.StatusInternalServerError, string(consts.ExprHttpInternalServerError))
			}
		}

		/*var cErr error
		if ctx.Request().Method == http.MethodHead {
			cErr = ctx.NoContent(code)
		} else {
			cErr = ctx.File(fmt.Sprintf("%d.html", code)) // eg: 404.html
		}
		if cErr != nil {
			ctx.Logger().Error("error page sending failed or not exists", "error", errors.Join(err, cErr))
		}*/
	}

	e.Use(middleware.Recover())

	addrport := net.JoinHostPort(cfg.Http.Address, fmt.Sprintf("%d", cfg.Http.Port))
	slog.Info("Listening on " + addrport)
	if err := e.Start(addrport); err != nil {
		slog.Error("FATAL: Failed to start HTTP server.")
		os.Exit(1)
	}

	/*
		// app
		app := New(cfg, db)
		app.echo.Logger = utils.GetLogger()
	*/
}
