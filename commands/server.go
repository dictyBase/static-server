package commands

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dictyBase/static-server/logger"
	"github.com/gorilla/handlers"
	"gopkg.in/urfave/cli.v1"
)

// rpm := regexp.MustCompile("precache-manifest")

func ServeAction(c *cli.Context) error {
	// create log folder
	lmw, err := logger.GetLoggerMiddleware(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	fs := handlers.CompressHandlerLevel(http.FileServer(http.Dir(c.String("folder"))), gzip.BestCompression)
	vhandler := fs
	port := fmt.Sprintf(":%d", c.Int("port"))
	subURL := fmt.Sprintf("/%s/", strings.TrimPrefix(c.String("static-folder"), "/"))
	sw := "service-worker.js"
	m := "manifest.json"
	fav := "favicon.ico"

	if len(c.String("sub-url")) > 0 {
		prefixPath := fmt.Sprintf("/%s", strings.TrimPrefix(c.String("sub-url"), "/"))
		subURL = fmt.Sprintf(
			"%s%s",
			prefixPath,
			subURL,
		)
		sw = fmt.Sprintf("%s/%s", c.String("sub-url"), sw)
		m = fmt.Sprintf("%s/%s", c.String("sub-url"), m)
		fav = fmt.Sprintf("%s/%s", c.String("sub-url"), fav)
		vhandler = http.StripPrefix(prefixPath, fs)
	}
	http.Handle(subURL, lmw.Middleware(vhandler))
	http.HandleFunc(sw, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("serving %s", sw)
		http.ServeFile(w, r, "service-worker.js")
	})
	http.HandleFunc(m, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("serving %s", m)
		http.ServeFile(w, r, "manifest.json")
	})
	http.HandleFunc(fav, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("serving %s", fav)
		http.ServeFile(w, r, "favicon.ico")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("serving client url %s", r.URL.Path)
		http.ServeFile(w, r, filepath.Join(c.String("folder"), "index.html"))
	})
	log.Printf("listening to port %s with url %s\n", port, subURL)
	log.Fatal(http.ListenAndServe(port, nil))
	return nil
}
