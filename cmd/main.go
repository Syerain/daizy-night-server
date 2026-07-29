package main

import (
	"daizynight/internal/config"
	"daizynight/internal/crypto"
	"daizynight/internal/db"
	"daizynight/internal/router"
	"daizynight/internal/utils"
	"fmt"
	"log"
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
	cfg := config.MustLoadConfig()

	// <--------- configuration below --------->

	// initializing our colored logger
	utils.InitModuleLogger(cfg.Main.IsDebugMode, "main")

	// initializing crypto
	err := crypto.Init(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	// initializing jwt
	crypto.InitJwt(cfg)

	// printing basic config info
	{
		slog.Info("Debug Mode:", slog.Bool("enabled", cfg.Main.IsDebugMode))
		slog.Info("Listening", "address", cfg.Http.Address, "port", cfg.Http.Port)
	}

	// <---------- service below --------->

	// starting db
	err = db.Init(cfg)
	if err != nil {
		log.Fatal("FATAL:Couldnt init database !")
	}

	// loading Echo engine
	slog.Info("Loading HTTP server ...")

	e := router.New()
	e.Logger = utils.GetLogger()
	addrport := net.JoinHostPort(cfg.Http.Address, fmt.Sprintf("%d", cfg.Http.Port))

	slog.Info("Listening on " + addrport)

	if err := e.Start(addrport); err != nil {
		slog.Error("FATAL: Failed to start HTTP server.")
		os.Exit(1)
	}
}
