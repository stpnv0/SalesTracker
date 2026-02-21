package main

import (
	"log"

	"github.com/stpnv0/SalesTracker/internal/app"
	"github.com/stpnv0/SalesTracker/internal/config"
	wbflogger "github.com/wb-go/wbf/logger"
)

func main() {
	cfg := config.MustLoad()

	logger, err := wbflogger.InitLogger(
		cfg.Logger.LogEngine(),
		"SalesTracker",
		cfg.Gin.Mode,
		wbflogger.WithLevel(cfg.Logger.LogLevel()),
	)
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}

	application, err := app.New(cfg, logger)
	if err != nil {
		log.Fatalf("app init: %v", err)
	}

	if err = application.Run(); err != nil {
		log.Fatalf("app run: %v", err)
	}
}
