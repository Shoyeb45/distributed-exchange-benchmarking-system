package main

import (
	"net/http"
	"time"

	_ "github.com/Shoyeb45/fast-docs/docs"
	"github.com/Shoyeb45/fast-docs/internal/app"
	"github.com/Shoyeb45/fast-docs/pkg/config"
	"github.com/Shoyeb45/fast-docs/pkg/database"
	"github.com/Shoyeb45/fast-docs/pkg/logger"
)

// @title           Server
// @version        	2.0
// @description     API documentation
// @host            localhost:8080
// @BasePath        /api/

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

	if err := srv.ListenAndServe(); err != nil {
		logger.Log.Error("failed to start application")
		panic(err.Error())
	}
}
