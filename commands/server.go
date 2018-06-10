package commands

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dictyBase/static-server/logger"
	"github.com/gorilla/handlers"
	"gopkg.in/urfave/cli.v1"
)

func ServeAction(c *cli.Context) error {
	// create log folder
	loggerMw, err := logger.GetLoggerMiddleware(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	fs := handlers.CompressHandlerLevel(http.FileServer(http.Dir(c.String("folder"))), gzip.BestCompression)
	port := fmt.Sprintf(":%d", c.Int("port"))
	if c.IsSet("sub-url") {
		subURL := c.String("sub-url")
		if !strings.HasSuffix(subURL, "/") {
			subURL = subURL + "/"
		}
		http.Handle(subURL, http.StripPrefix(subURL, loggerMw.Middleware(fs)))
		log.Printf("listening to port %s with url %s\n", port, subURL)
	} else {
		http.Handle("/", loggerMw.Middleware(fs))
		log.Printf("listening to port %s\n", port)
	}
	log.Fatal(http.ListenAndServe(port, nil))
	return nil
}
