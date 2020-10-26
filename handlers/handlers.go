package handlers

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/urfave/cli"
)

type AssetHandler struct {
	File string
}

func (a *AssetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, a.File)
}

type RootHandler struct {
	PreRegexp      *regexp.Regexp
	Folder, SubURL string
}

func NewRootHandler(url, f string) *RootHandler {
	return &RootHandler{
		SubURL:    url,
		Folder:    f,
		PreRegexp: regexp.MustCompile(`(?m)precache-manifest\.\w+\.js`),
	}
}

func (rh *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rh.PreRegexp.FindString(r.URL.Path) != "" {
		url := strings.TrimPrefix(r.URL.Path, "/")
		if len(rh.SubURL) > 0 {
			url = strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%s%s", rh.SubURL, "/"))
		}
		http.ServeFile(w, r, url)
		return
	}
	http.ServeFile(w, r, filepath.Join(rh.Folder, "index.html"))
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
