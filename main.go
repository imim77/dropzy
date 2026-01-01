package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	cfg := &Config{Port: ":42069"}
	logger := NewLogger(os.Stdout)
	srv := NewServer(cfg, logger)
	httpServer := &http.Server{
		Addr:    cfg.Port,
		Handler: srv,
	}

	go func() {
		logger.InfoMess("Server starting on: " + httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Error("http server failed", err)
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	return httpServer.Shutdown(shutdownCtx)
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
