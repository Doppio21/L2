package main

import (
	"L2/develop/dev11/server"
	"context"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	address := "localhost:10001"
	events := server.NewEvents()
	log := logrus.New()
	cfg := server.Config{
		Address: address,
	}
	deps := server.Dependencies{
		Events: events,
		Log:    log,
	}
	srv := server.New(cfg, deps)
	srv.Run(ctx)
}
