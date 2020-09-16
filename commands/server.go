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
	"github.com/urfave/cli"
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
	swPath := fmt.Sprintf("/%s", sworker)
	manifestPath := fmt.Sprintf("/%s", manifest)
	faviconPath := fmt.Sprintf("/%s", favicon)

	if len(c.String("sub-url")) > 0 {
		prefixPath := fmt.Sprintf("/%s", strings.TrimPrefix(c.String("sub-url"), "/"))
		subURL = fmt.Sprintf(
			"%s%s",
			prefixPath,
			subURL,
		)
		swPath = fmt.Sprintf("%s/%s", c.String("sub-url"), sworker)
		manifestPath = fmt.Sprintf("%s/%s", c.String("sub-url"), manifest)
		faviconPath = fmt.Sprintf("%s/%s", c.String("sub-url"), favicon)
		vhandler = http.StripPrefix(prefixPath, fs)
	}
	http.Handle(subURL, lmw.Middleware(vhandler))
	http.HandleFunc(swPath, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("serving service-worker file %s", swPath)
		http.ServeFile(w, r, sworker)
	})
	http.HandleFunc(manifestPath, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("serving manifest.json %s", manifestPath)
		http.ServeFile(w, r, manifest)
	})
	http.HandleFunc(faviconPath, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("serving favicon file %s", faviconPath)
		http.ServeFile(w, r, favicon)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if precacheRegex.FindString(r.URL.Path) != "" {
			url := strings.TrimPrefix(r.URL.Path, "/")
			if len(c.String("sub-url")) > 0 {
				url = strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%s%s", c.String("sub-url"), "/"))
			}
			log.Printf("serving precache-manifest file %s", url)
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
