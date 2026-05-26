package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/Shoyeb45/server/pkg/config"
)

// Log is Logger and it should be used througout in the application.
var Log *slog.Logger

// Init will Initialize the logger, must be called before accessing `Log`.
func Init() error {
	stage := config.Cfg.Stage

	var handler slog.Handler

	// in production, we should have structured logging inside the file and console and in dev
	// we should have in the stdout
	switch stage {
	case "prod":
		// first make the directory
		if err := os.MkdirAll(config.Cfg.LogDirectory, 0750); err != nil {
			panic(err.Error())
		}

		file, err := os.OpenFile("./"+config.Cfg.LogDirectory+"/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)

		if err != nil {
			return err
		}
		// log both in the file and in stdout
		mutiwriter := io.MultiWriter(os.Stdout, file)

		handler = slog.NewJSONHandler(mutiwriter, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})

	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}
	Log = slog.New(handler).With(
		"service", "server",
		"env", stage,
	)
	slog.SetDefault(Log)

	return nil
}
