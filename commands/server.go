package commands

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"

	"github.com/dictyBase/go-middlewares/middlewares/cache"
	"github.com/dictyBase/go-middlewares/middlewares/chain"
	"github.com/dictyBase/go-middlewares/middlewares/nocache"
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
	cacheMw := cache.NewHTTPCache(c.Int("cache-duration"))
	cacheChain := chain.NewChain(lmw.Middleware, cacheMw.Middleware)
	nocacheChain := chain.NewChain(lmw.Middleware, nocache.Middleware)
	subURL, baseHandler := sh.BaseHandlerAndPath(c)
	mux := http.NewServeMux()
	mux.Handle(subURL, cacheChain.Then(baseHandler))
	for file, path := range pathMap(c.String("sub-url")) {
		mux.Handle(path, cacheChain.Then(
			handlers.CompressHandlerLevel(
				&sh.AssetHandler{File: file},
				gzip.BestCompression,
			),
		))
	}
	rh := sh.NewRootHandler(subURL, c.String("folder"))
	mux.Handle("/", nocacheChain.Then(
		handlers.CompressHandlerLevel(
			rh,
			gzip.BestCompression,
		),
	))
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
