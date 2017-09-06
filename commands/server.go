package commands

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dictyBase/static-server/logger"
	"gopkg.in/urfave/cli.v1"
)

func ServeAction(c *cli.Context) error {
	// create log folder
	loggerMw, err := logger.GetLoggerMiddleware(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	fs := http.FileServer(http.Dir(c.String("folder")))
	http.Handle("/", loggerMw.Middleware(fs))
	port := fmt.Sprintf(":%d", c.Int("port"))
	log.Printf("listening to port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
	return nil
}
