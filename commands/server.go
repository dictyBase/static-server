package commands

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dictyBase/static-server/logger"
	"github.com/gorilla/handlers"
	"gopkg.in/urfave/cli.v1"
)

func sendSW(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("service-worker.js")
	if err != nil {
		http.Error(w, "Couldn't read service worker file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Write(data)
}

func sendManifest(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("manifest.json")
	if err != nil {
		http.Error(w, "Couldn't read manifest.json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

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
	if len(c.String("sub-url")) > 0 {
		prefixPath := fmt.Sprintf("/%s", strings.TrimPrefix(c.String("sub-url"), "/"))
		subURL = fmt.Sprintf(
			"%s%s",
			prefixPath,
			subURL,
		)
		vhandler = http.StripPrefix(prefixPath, fs)
	}
	http.Handle(subURL, lmw.Middleware(vhandler))
	http.HandleFunc("/service-worker.js", sendSW)
	http.HandleFunc("/manifest.json", sendManifest)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("serving client url %s", r.URL.Path)
		http.ServeFile(w, r, filepath.Join(c.String("folder"), "index.html"))
	})
	log.Printf("listening to port %s with url %s\n", port, subURL)
	log.Fatal(http.ListenAndServe(port, nil))
	return nil
}
