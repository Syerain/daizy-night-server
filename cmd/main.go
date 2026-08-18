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
	cfg := config.MustLoadConfig()

	srv, err := New(cfg)
	if err != nil {
		slog.Error("FATAL: failed to init server", slog.Any("error", err))
		os.Exit(1)
	}

	if err := srv.Run(); err != nil {
		slog.Error("FATAL: server exited with error", slog.Any("error", err))
		os.Exit(1)
	}
}

type Server struct {
	e   *echo.Echo
	cfg *config.Config
	pDB abstract.InterfaceProviderDB
}

// New assembles every dependency (logger, crypto, database, repos, services,
// handlers and the echo engine) and returns a fully-initialized server.
// It never calls os.Exit: all failures are returned to the caller.
func New(cfg *config.Config) (*Server, error) {
	utils.InitModuleLogger(cfg.Main.IsDebugMode, "main")
	slog.Info("Server init..")

	// config
	slog.Info("Debug Mode:", slog.Bool("enabled", cfg.Main.IsDebugMode))
	slog.Info("Configured listening at", "address", cfg.Http.Address, "port", cfg.Http.Port)

	// crypto provider
	pCrypto, err := crypto.NewProviderCrypto(cfg)
	if err != nil {
		return nil, fmt.Errorf("init crypto: %w", err)
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
		return nil, fmt.Errorf("init db: %w", err)
	}

	// dbware repo
	repoUser := dbware.NewRepoUser(pDB)
	repoToken := dbware.NewRepoToken(pDB, cfg)
	repoRegcode := dbware.NewRepoRegistercode(pDB)

	// service provider
	pSvcCode := service.NewServiceCode(repoRegcode)
	pSvcUser := service.NewServiceUser(repoUser, repoToken, pSvcCode, pCrypto)
	pSvcAdmin := service.NewServiceAdmin()
	s := service.NewServiceComplex(pSvcUser, pSvcCode, pSvcAdmin)

	// handler
	h := handler.NewHandlerComplex(s.ServiceUser, s.ServiceCode, s.ServiceAdmin)

	// Echo router
	e := router.New(h, pCrypto, cfg)
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

	return &Server{
		e:   e,
		cfg: cfg,
		pDB: pDB,
	}, nil
}

// Run starts the HTTP server and blocks until it is shut down. echo v5 handles
// SIGINT/SIGTERM with a graceful shutdown internally, so a nil return means the
// server stopped on purpose. The database is closed right before returning.
func (s *Server) Run() error {
	defer func() {
		if err := s.pDB.Close(); err != nil {
			slog.Warn("failed to close database", slog.Any("error", err))
		}
	}()

	addrport := net.JoinHostPort(s.cfg.Http.Address, fmt.Sprintf("%d", s.cfg.Http.Port))
	slog.Info("Listening on " + addrport)

	err := s.e.Start(addrport)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	slog.Info("Server shut down gracefully")
	return nil
}

func formatLine(
	chain []string,
	method string,
	path string,
	err error,
) string {
	l := fmt.Sprintf(
		"Failure-Request\tmethod::%s\terror::%s\tcallchain::%s\tpath::%s\tdetails::%s;",
		fmt.Sprint(chain), method, fmt.Sprintf("%T", err), path, err.Error(),
	)
	return l
}
