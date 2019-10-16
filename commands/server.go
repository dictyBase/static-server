package commands

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dictyBase/static-server/logger"
	"github.com/gorilla/handlers"
	"gopkg.in/urfave/cli.v1"
)

const (
	sworker  = "service-worker.js"
	manifest = "manifest.json"
	favicon  = "favicon.ico"
)

var precacheRegex = regexp.MustCompile(`(?m)precache-manifest\.\w+\.js`)

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
	swPath := sw
	manifestPath := manifest
	faviconPath := favicon

	if len(c.String("sub-url")) > 0 {
		prefixPath := fmt.Sprintf("/%s", strings.TrimPrefix(c.String("sub-url"), "/"))
		subURL = fmt.Sprintf(
			"%s%s",
			prefixPath,
			subURL,
		)
		swPath = fmt.Sprintf("%s/%s", c.String("sub-url"), sw)
		manifestPath = fmt.Sprintf("%s/%s", c.String("sub-url"), m)
		faviconPath = fmt.Sprintf("%s/%s", c.String("sub-url"), fav)
		vhandler = http.StripPrefix(prefixPath, fs)
	}
	http.Handle(subURL, lmw.Middleware(vhandler))
	http.HandleFunc(swPath, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, sworker)
	})
	http.HandleFunc(manifestPath, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, manifest)
	})
	http.HandleFunc(faviconPath, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, favicon)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if precacheRegex.FindString(r.URL.Path) != "" {
			var url string
			if len(c.String("sub-url")) > 0 {
				url = strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%s%s", c.String("sub-url"), "/"))
			}
			log.Printf("serving precache-manifest %s", url)
			http.ServeFile(w, r, url)
			return
		}
		log.Printf("serving client url %s", r.URL.Path)
		http.ServeFile(w, r, filepath.Join(c.String("folder"), "index.html"))
	})
	log.Printf("listening to port %s with url %s\n", port, subURL)
	log.Fatal(http.ListenAndServe(port, nil))
	return nil
}
