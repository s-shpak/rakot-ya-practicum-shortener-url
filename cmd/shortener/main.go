package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rakot9/ya-practicum-shortener-url/internal/config"
	"github.com/rakot9/ya-practicum-shortener-url/internal/router"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	flags := config.ParseFlags(os.Args[1:])

	router := router.Router(flags.FlagRunAddr, flags.FlagRunShorternerAddr, flags.FlagLog)

	if err := run(router, flags); err != nil {
		slog.Error("failed to load configuration. Error: ", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(router chi.Router, flags config.Flags) error {

	if flags.FlagLog {
		slog.Info("Running server on", slog.Any("flag_run_addr", flags.FlagRunAddr))
	}

	return http.ListenAndServe(flags.FlagRunAddr, router)
}
