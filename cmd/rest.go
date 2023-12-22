package main

import (
	"context"
	"errors"
	"golinkcut/api/rest"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func runRestApi(uc link.UseCase, cfg config.Config) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	r := rest.SetupRouter(uc, cfg)
	srv := &http.Server{
		Addr:    ":" + cfg["httpPort"].(string),
		Handler: r,
	}
	log.Printf("REST API started at port %s", cfg["httpPort"].(string))
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server listen error: %s\n", err)
		}
	}()
	// Listen for the interrupt signal.
	<-ctx.Done()
	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down HTTP server gracefully, press Ctrl+C again to force")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server forced to shutdown: %s\n", err)
	}
	log.Println("HTTP server exiting")
}
