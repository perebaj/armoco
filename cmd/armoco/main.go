// armoco runs the service.
// For details on how to configure it just run:
//
//	armoco --help
package main

import (
	"os"

	"github.com/birdie-ai/golibs/slog"
)

func main() {
	logcfg, err := slog.LoadConfig("ARMOCO")
	abortonerr(err)

	err = slog.Configure(logcfg)
	abortonerr(err)

	slog.Info("TODO: implement armoco")
}

func abortonerr(err error) {
	if err != nil {
		slog.Error("fatal error initializing service", "error", err.Error())
		os.Exit(1)
	}
}
