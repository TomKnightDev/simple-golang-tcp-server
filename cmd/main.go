package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/tomknightdev/tcp/server"
	"github.com/tomknightdev/tcp/store"
)

func main() {
	logger := log.New(os.Stdout, "tcp-server", log.LstdFlags|log.Lshortfile)
	store := store.New(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	server.StartTCPServer(ctx, logger, *store)

	logger.Println("server stopped")
}
