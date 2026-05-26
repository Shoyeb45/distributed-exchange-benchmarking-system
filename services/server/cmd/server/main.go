package main

import (
	"net/http"
	"time"

	"github.com/Shoyeb45/server/api/modules"
	_ "github.com/Shoyeb45/server/docs"
	"github.com/Shoyeb45/server/internal/app"
	"github.com/Shoyeb45/server/pkg/config"
	"github.com/Shoyeb45/server/pkg/database"
	"github.com/Shoyeb45/server/pkg/logger"
)

// @title           Server
// @version        	2.0
// @description     API documentation
// @host            localhost:8000
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// read environment variables
	if err := config.LoadEnvironmentVariables(); err != nil {
		panic(err.Error())
	}

	// initialize logger
	if err := logger.Init(); err != nil {
		panic(err.Error())
	}

	logger.Log.Info("Logger initialized successfully.")

	if err := database.Connect(); err != nil {
		panic(err.Error())
	}
	defer database.Close()

	chiMux := app.New()

	const readAndWriteTimeout = 10 * time.Second
	const idleTimeout = 60 * time.Second

	srv := &http.Server{
		Addr:         ":" + config.Cfg.Port,
		Handler:      chiMux,
		ReadTimeout:  readAndWriteTimeout,
		WriteTimeout: readAndWriteTimeout,
		IdleTimeout:  idleTimeout,
	}

	// mount all the routes
	modules.MountRoutes(chiMux);

	if err := srv.ListenAndServe(); err != nil {
		logger.Log.Error("failed to start application")
		panic(err.Error())
	}
}
