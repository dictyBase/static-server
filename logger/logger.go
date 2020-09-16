package logger

import (
	"fmt"
	"io"
	"os"

	loggerMw "github.com/dictyBase/go-middlewares/middlewares/logrus"

	"github.com/urfave/cli"
)

// GetLoggerMiddleware gets a net/http compatible instance of logrus
func GetLoggerMiddleware(c *cli.Context) (*loggerMw.Logger, error) {
	var logger *loggerMw.Logger
	var w io.Writer
	if c.IsSet("log-file") {
		fw, err := os.Create(c.String("log-file"))
		if err != nil {
			return logger,
				fmt.Errorf("could not open log file  %s %s", c.String("log-file"), err)
		}
		w = io.MultiWriter(fw, os.Stderr)
	} else {
		w = os.Stderr
	}
	if c.String("log-format") == "json" {
		logger = loggerMw.NewJSONFileLogger(w)
	} else {
		logger = loggerMw.NewFileLogger(w)
	}
	return logger, nil
}
