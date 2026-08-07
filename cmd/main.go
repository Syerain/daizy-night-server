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

// entrance
func main() {
	Run()
}

type Server struct {
	echo *echo.Echo
	cfg  *config.Config
	db   *dbware.ProviderDB
}

// havent use
func NewServer(cfg *config.Config, dbProvider *dbware.ProviderDB) *Server {
	// logger
	utils.InitModuleLogger(cfg.Main.IsDebugMode, "main")

	return &Server{
		echo: echo.New(),
		cfg:  cfg,
		db:   dbProvider,
	}
}

func Run() {
	// start
	slog.Info("Server starting...")
	slog.Info("Reading config...")

	// config
	cfg := config.MustLoadConfig()
	utils.InitModuleLogger(cfg.Main.IsDebugMode, "main")
	slog.Info("Debug Mode:", slog.Bool("enabled", cfg.Main.IsDebugMode))
	slog.Info("Configured listening at", "address", cfg.Http.Address, "port", cfg.Http.Port)

	// crypto provider
	pCrypto, err := crypto.NewProviderCrypto(cfg)
	if err != nil {
		slog.Error("FATAL: couldnt init crypto !")
		os.Exit(1)
	}

	// database provider
	pDB, err := dbware.NewProviderDB(struct {
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
	defer pDB.Close()

	// dbware repo
	repoUser := dbware.NewRepoUser(pDB)
	repoToken := dbware.NewRepoToken(pDB)
	repoRegcode := dbware.NewRepoRegistercode(pDB)

	// service provider
	svcUser := service.NewServiceUser(repoUser, repoToken, repoRegcode, pCrypto)
	svcCode := service.NewServiceCode(repoRegcode)
	svcAdmin := service.NewServiceAdmin()

	// handler
	handler := &handler.HandlerComplex{
		ServiceUser:  svcUser,
		ServiceCode:  svcCode,
		ServiceAdmin: svcAdmin,
	}

	// Echo router
	e := router.New(handler, pCrypto, cfg)
	// set Echo logger
	e.Logger = utils.GetLogger()
	// error handler
	e.HTTPErrorHandler = func(ctx *echo.Context, err error) {
		// return if committed
		if resp, uErr := echo.UnwrapResponse(ctx.Response()); uErr == nil {
			if resp.Committed {
				return
			}
		}

		// logger layer;
		utils.Layer(ctx).Error(formatLine(
			utils.GetCallChain(ctx),
			ctx.Request().Method,
			ctx.Request().URL.Path,
			err,
		))
		/*utils.Layer(ctx).Error("request failure",
			slog.String("callchain", fmt.Sprint(utils.GetCallChain(ctx))),
			slog.String("method", ctx.Request().Method),
			slog.String("path", ctx.Request().URL.Path),
			slog.String("errtype", fmt.Sprintf("%T", err)),
			slog.Any("details", err.Error()),
		)*/

		// custom errors implements func Respond() and mid.RespondCustom() calls it to echo User-Friendly msg.
		if errapp, ok := errs.Easx[abstract.InterfaceAppError](err); ok {
			mid.RespondCustom(ctx, errapp)
			return
		}

		// http status code
		var sc echo.HTTPStatusCoder
		stat := http.StatusInternalServerError // set defult status code to IDE as a fallback
		if errors.As(err, &sc) {
			if tmp := sc.StatusCode(); tmp != 0 {
				stat = tmp
			}
		}

		// final response
		mid.Respond(ctx, stat, string(consts.HttpExprInternalServerError))
	}

	// using Echo default Recover()
	e.Use(middleware.Recover())

	// run http server
	addrport := net.JoinHostPort(cfg.Http.Address, fmt.Sprintf("%d", cfg.Http.Port))
	slog.Info("Listening on " + addrport)
	if err := e.Start(addrport); err != nil {
		slog.Error("FATAL: Failed to start HTTP server.")
		os.Exit(1)
	}
}

func formatLine(
	chain []string,
	method string,
	path string,
	err error,
) string {
	l := fmt.Sprintf(
		"Failure-Request\tcallchain::%s\tmethod::%s\tpath::%s\tdetails::%s",
		fmt.Sprint(chain), method, path, fmt.Sprintf("%T", err),
	)
	return l
}
