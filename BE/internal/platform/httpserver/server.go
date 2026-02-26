package httpserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type ServerConfig struct {
	Name            string
	Port            string
	Handler         http.Handler
	ShutdownTimeout time.Duration
}

func Run(ctx context.Context, cfg ServerConfig) error {
	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = 10 * time.Second
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           cfg.Handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("[%s] http listening on :%s", cfg.Name, cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		log.Printf("[%s] shutting down http server", cfg.Name)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
