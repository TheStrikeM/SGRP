package main

import (
	"fmt"
	"log/slog"
	"os"
	"sgrp/internal/lib/logger/handlers/slogpretty"
	sgrp "sgrp/internal/protocol"
	"sync"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type LoginForm struct {
	login    string
	password string
}

func main() {
	log := setupLogger(envLocal)

	var wg sync.WaitGroup
	wg.Add(1)

	sgrpServer := sgrp.New(log, 8082)
	sgrpServer.AddRoute("/TEST", func(request sgrp.StrikeRequest) sgrp.StrikeResponse {
		return sgrp.StrikeResponse{
			Result: request.Body + " tested test",
		}
	})
	sgrpServer.AddRoute("/HELLO", func(request sgrp.StrikeRequest) sgrp.StrikeResponse {
		return sgrp.StrikeResponse{
			Result: fmt.Sprintf("%s tested hello", request.Body),
		}
	})

	waitGroup := sgrpServer.MustRun()
	waitGroup.Wait()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
