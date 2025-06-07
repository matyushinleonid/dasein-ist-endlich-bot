package listener

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
)

func serveHealth(ctx context.Context, logger logr.Logger, port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	go func() {
		logger.Info("starting health server", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(err, "health server crashed")
		}
	}()
	go func() {
		<-ctx.Done()
		_ = srv.Close()
	}()
}
