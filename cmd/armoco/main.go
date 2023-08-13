// armoco runs the service.
// For details on how to configure it just run:
//
//	armoco --help
package main

import (
	"net/http"
	"os"
	"time"

	"github.com/birdie-ai/armoco/api"
	"github.com/birdie-ai/golibs/slog"
)

func main() {
	logcfg, err := slog.LoadConfig("ARMOCO")
	abortonerr(err)

	err = slog.Configure(logcfg)
	abortonerr(err)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      api.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	slog.Info("starting service", "addr", srv.Addr)
	err = srv.ListenAndServe()
	abortonerr(err)

	slog.Info("TODO: implement armoco")
}

func abortonerr(err error) {
	if err != nil {
		slog.Error("fatal error initializing service", "error", err.Error())
		os.Exit(1)
	}
}
