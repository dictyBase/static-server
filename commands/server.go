package commands

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"

	sh "github.com/dictyBase/static-server/handlers"
	"github.com/dictyBase/static-server/logger"
	"github.com/gorilla/handlers"
	"github.com/urfave/cli"
)

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

func ServeAction(c *cli.Context) error {
	lmw, err := logger.GetLoggerMiddleware(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	subURL, baseHandler := sh.BaseHandlerAndPath(c)
	mux := http.NewServeMux()
	mux.Handle(subURL, lmw.Middleware(baseHandler))
	for file, path := range pathMap(c.String("sub-url")) {
		mux.Handle(path, lmw.Middleware(
			handlers.CompressHandlerLevel(
				&sh.AssetHandler{File: file},
				gzip.BestCompression,
			),
		))
	}
	rh := sh.NewRootHandler(subURL, c.String("folder"))
	mux.Handle("/", hanlders.NoCache(
		lmw.Middleware(
			handlers.CompressHandlerLevel(
				rh,
				gzip.BestCompression,
			),
		)))
	port := fmt.Sprintf(":%d", c.Int("port"))
	log.Printf("listening to port %s with url %s\n", port, subURL)
	if err := http.ListenAndServe(port, mux); err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in running server %s", err),
			2,
		)
	}
	return nil
}
