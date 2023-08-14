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
	"github.com/cloudflare/cloudflare-go"
	"github.com/kelseyhightower/envconfig"
)

const serviceName = "armoco"

type Config struct {
	LogLevel           string `envconfig:"LOG_LEVEL"`
	LogFormat          string `envconfig:"LOG_FMT"`
	CloudFlareAPIToken string `envconfig:"CLOUDFLARE_API_TOKEN"`
}

func main() {
	logcfg, err := slog.LoadConfig(serviceName)
	abortonerr(err)
	var cfg Config
	if err := envconfig.Process(serviceName, &cfg); err != nil {
		slog.Fatal("failed to load config from environment", "error", err.Error())
	}

	err = slog.Configure(logcfg)

	cloudFlareClient, err := openCloudFlare(cfg)
	if err != nil {
		slog.Fatal("failed to open cloudflare client", "error", err.Error())
	}

	app := api.Application{
		CloudFlare: cloudFlareClient,
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
}

func abortonerr(err error) {
	if err != nil {
		slog.Error("fatal error initializing service", "error", err.Error())
		os.Exit(1)
	}
}

func openCloudFlare(cfg Config) (cloudFlareAPI *cloudflare.API, err error) {
	api, err := cloudflare.NewWithAPIToken(cfg.CloudFlareAPIToken)
	if err != nil {
		return nil, err
	}

	return api, nil
}
