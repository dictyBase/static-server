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

var precacheRegex = regexp.MustCompile(`(?m)precache-manifest\.\w+\.js`)

func defaultStaticAssets() []string {
	return []string{
		"service-worker.js",
		"manifest.json",
		"favicon.ico",
	}
}

func pathMap(subURL string) map[string]string {
	pm := make(map[string]string)
	if len(subURL) > 0 {
		for _, p := range defaultStaticAssets() {
			pm[p] = fmt.Sprintf("%s/%s", subURL, p)
		}
	} else {
		for _, p := range defaultStaticAssets() {
			pm[p] = fmt.Sprintf("/%s", p)
		}
	}
	return pm
}

func BaseHandlerAndPath(c *cli.Context) (string, http.Handler) {
	h := handlers.CompressHandlerLevel(
		http.FileServer(
			http.Dir(c.String("folder")),
		),
		gzip.BestCompression,
	)
	subURL := fmt.Sprintf("/%s/", strings.TrimPrefix(c.String("static-folder"), "/"))
	if len(c.String("sub-url")) > 0 {
		prefixPath := fmt.Sprintf("/%s", strings.TrimPrefix(c.String("sub-url"), "/"))
		h = http.StripPrefix(prefixPath, h)
		subURL = fmt.Sprintf(
			"%s%s",
			prefixPath,
			subURL,
		)
	}
	return subURL, h
}

func ServeAction(c *cli.Context) error {
	// create log folder
	lmw, err := logger.GetLoggerMiddleware(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	fs := handlers.CompressHandlerLevel(http.FileServer(http.Dir(c.String("folder"))), gzip.BestCompression)
	vhandler := fs
	subURL := fmt.Sprintf("/%s/", strings.TrimPrefix(c.String("static-folder"), "/"))
	if len(c.String("sub-url")) > 0 {
		prefixPath := fmt.Sprintf("/%s", strings.TrimPrefix(c.String("sub-url"), "/"))
		subURL = fmt.Sprintf(
			"%s%s",
			prefixPath,
			subURL,
		)
		vhandler = http.StripPrefix(prefixPath, fs)
		for file, path := range customPathMap(c.String("sub-url")) {
			http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, file)
			})
		}
	} else {
		for file, path := range defaultPathMap() {
			http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, file)
			})
		}

	}
	http.Handle(subURL, lmw.Middleware(vhandler))
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
	port := fmt.Sprintf(":%d", c.Int("port"))
	log.Printf("listening to port %s with url %s\n", port, subURL)
	log.Fatal(http.ListenAndServe(port, nil))
	return nil
}
