package prometheus

import (
	"context"
	"net/http"
	"sync"

	"github.com/deepshore/chirpstack-simulator/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func Setup(ctx context.Context, wg *sync.WaitGroup, bind string) error {
	log.Info("Setting up Prometheus metrics endpoint")

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Handler: mux,
		Addr:    bind,
	}

	log.Info("Prometheus metrics server starting on ", config.C.Prometheus.Bind)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start Prometheus metrics server")
		}
	}()

	return nil
}
