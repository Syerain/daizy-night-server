package main

import (
	"daizynight/internal/config"
	"daizynight/internal/crypto"
	"daizynight/internal/dbware"
	"daizynight/internal/handler"
	"daizynight/internal/router"
	"daizynight/internal/service"
	"daizynight/internal/utils"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/labstack/echo/v5"
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

	// router
	e := router.New(handler)
	e.Logger = utils.GetLogger()
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
