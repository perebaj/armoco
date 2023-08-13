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
	"github.com/kelseyhightower/envconfig"
)

const serviceName = "armoco"

type cfg struct {
	CfApiKey   string `envconfig:"CF_API_KEY"`
	CfApiEmail string `envconfig:"CF_API_EMAIL"`
}

func main() {
	logcfg, err := slog.LoadConfig(serviceName)
	abortonerr(err)

	if err := envconfig.Process(serviceName, &cfg{}); err != nil {
		slog.Fatal("failed to load config from environment", "error", err.Error())
	}

	err = slog.Configure(logcfg)
	abortonerr(err)

	// _, err = cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))

	// abortonerr(err)

	app := &api.Application{
		Cloudflare: nil,
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      app.Routes(),
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
