package main

import (
	"daizynight/internal/config"
	"daizynight/internal/crypto"
	"daizynight/internal/db"
	"daizynight/internal/router"
	"daizynight/internal/utils"
	"fmt"
	"log/slog"
	"net"
	"os"
)

func main() {
	// reading config, temply using slog vanilla logger.
	{
		slog.Info("Server starting...")
		slog.Info("Reading config...")
	}

	// initializing Global Config
	err := config.InitConfig()
	if err != nil {
		slog.Error("FATAL: Failed to initialize config", "error", err)
		return
	}

	cfg := config.GetConfig()

	// <--------- configuration below --------->

	// initializing our colored logger
	utils.InitModuleLogger(cfg.Main.IsDebugMode, "main")

	// initializing crypto
	err = crypto.Init(cfg)
	if err != nil {
		slog.Error("Failed to init Crypto module !")
	}

	// initializing jwt
	crypto.InitJwt(cfg)

	// printing basic config info
	{
		slog.Info("Debug Mode:", slog.Bool("debug", cfg.Main.IsDebugMode))
		slog.Info("Config listen port: " + fmt.Sprintf("%d", cfg.Http.ListenPort))
		slog.Info("Config listen address: " + cfg.Http.ListenAddress)
	}

	// <---------- service below --------->

	// starting db
	err = db.Init(cfg)
	if err != nil {
		slog.Error("FATAL:Couldnt init database !")
	}

	// loading Echo engine
	slog.Info("Loading HTTP server ...")

	e := router.New()
	e.Logger = utils.GetLogger()
	addrport := net.JoinHostPort(cfg.Http.ListenAddress, fmt.Sprintf("%d", cfg.Http.ListenPort))

	slog.Info("Listening on " + addrport)

	if err := e.Start(addrport); err != nil {
		slog.Error("FATAL: Failed to start HTTP server.")
		os.Exit(1)
	}
}
