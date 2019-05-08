package netkat

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"io"
)

var Logger log.Logger

// InitJsonLogger function initiates a structured JSON logger, taking in the specified log level for what is displayed at runtime.
func InitLogger(writer io.Writer, logLevel string) {
	Logger = log.NewLogfmtLogger(writer)
	Logger = log.With(Logger, "timestamp", log.DefaultTimestampUTC)
	switch logLevel {
	case "debug":
		Logger = level.NewFilter(Logger, level.AllowDebug())
	case "info":
		Logger = level.NewFilter(Logger, level.AllowInfo())
	case "error":
		Logger = level.NewFilter(Logger, level.AllowError())
	default:
		Logger = level.NewFilter(Logger, level.AllowInfo())
	}
}
