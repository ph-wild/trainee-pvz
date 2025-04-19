package metrics

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"log/slog"
	"net/http"
	"time"
)

func RunMetricServer(listenAddress string) {
	mh := chi.NewRouter()
	mh.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	srv := &http.Server{
		Addr:        fmt.Sprintf(":%s", listenAddress),
		Handler:     mh,
		ReadTimeout: 1 * time.Second,
	}

	slog.Info("starting Metric exporter server: listening on", slog.String("address", listenAddress))

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		slog.Error("failed to listen promhandler server", "error", err)
	}
}
