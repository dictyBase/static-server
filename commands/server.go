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

var statRegxp = regexp.MustCompile(`.+(jpg|png|css|js)$`)

func ServeAction(c *cli.Context) error {
	// create log folder
	lmw, err := logger.GetLoggerMiddleware(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	fs := handlers.CompressHandlerLevel(http.FileServer(http.Dir(c.String("folder"))), gzip.BestCompression)
	idx := filepath.Join(c.String("folder"), "index.html")
	port := fmt.Sprintf(":%d", c.Int("port"))
	if c.IsSet("sub-url") {
		subURL := c.String("sub-url")
		if !strings.HasSuffix(subURL, "/") {
			subURL = subURL + "/"
		}
		http.Handle(subURL, lmw.Middleware(serveFromSubURL(subURL, idx, fs)))
		log.Printf("listening to port %s with url %s\n", port, subURL)
	} else {
		http.Handle("/", lmw.Middleware(serveFromRoot(idx, fs)))
		log.Printf("listening to port %s\n", port)
	}
	log.Fatal(http.ListenAndServe(port, nil))
	return nil
}

func serveFromRoot(idx string, fs http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if statRegxp.MatchString(r.URL.Path) {
			fs.ServeHTTP(w, r)
		} else {
			http.ServeFile(w, r, idx)
		}
	}
	return http.HandlerFunc(fn)
}

func serveFromSubURL(subURL string, idx string, fs http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if statRegxp.MatchString(r.URL.Path) {
			http.StripPrefix(subURL, fs).ServeHTTP(w, r)
		} else {
			http.ServeFile(w, r, idx)
		}
	}
	return http.HandlerFunc(fn)
}
