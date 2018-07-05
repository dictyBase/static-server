package commands

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/dictyBase/static-server/logger"
	"github.com/gorilla/handlers"
	"gopkg.in/urfave/cli.v1"
)

func ServeAction(c *cli.Context) error {
	// create log folder
	lmw, err := logger.GetLoggerMiddleware(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	fs := handlers.CompressHandlerLevel(http.FileServer(http.Dir(c.String("folder"))), gzip.BestCompression)
	port := fmt.Sprintf(":%d", c.Int("port"))
	vfolder := c.String("folder")
	if len(c.String("virtual-static-folder")) > 0 {
		vfolder = c.String("virtual-static-folder")
	}
	subURL := vfolder + "/"
	vhandler := http.StripPrefix(subURL, fs)
	if len(c.String("sub-url")) > 0 {
		subURL = fmt.Sprintf("%s%s/", c.String("sub-url"), vfolder)
		vhandler = http.StripPrefix(c.String("sub-url"), fs)
	}
	http.Handle(subURL, lmw.Middleware(vhandler))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("serving client url %s", r.URL.Path)
		http.ServeFile(w, r, filepath.Join(c.String("folder"), "index.html"))
	})
	log.Printf("listening to port %s with url %s\n", port, subURL)
	log.Fatal(http.ListenAndServe(port, nil))
	return nil
}
