package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/deepshore/chirpstack-simulator/internal/as"
	"github.com/deepshore/chirpstack-simulator/internal/config"
	"github.com/deepshore/chirpstack-simulator/internal/ns"
	"github.com/deepshore/chirpstack-simulator/internal/prometheus"
	"github.com/deepshore/chirpstack-simulator/internal/simulator"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) error {
	setLogLevel()
	tasks := []func(ctx context.Context, wg *sync.WaitGroup) error{
		setupApplicationClientAPIServer,
		setupNetworkServer,
		setupPrometheus,
		startSimulator,
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, task := range tasks {
		if err := task(ctx, &wg); err != nil {
			log.Fatalf("Failed to setup task: %v", err)
			return err
		}
	}

	exitChan := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	go func() {
		cancel()
		wg.Wait()
		exitChan <- struct{}{}
	}()
	cancel()

	select {
	case <-exitChan:
	case s := <-sigChan:
		log.Infof("Received signal: %s, shutting down gracefully", s)
	}

	return nil
}

func setupApplicationClientAPIServer(ctx context.Context, wg *sync.WaitGroup) error {
	return as.Setup(config.C)
}

func setupNetworkServer(ctx context.Context, wg *sync.WaitGroup) error {
	return ns.Setup(config.C)
}

func setupPrometheus(ctx context.Context, wg *sync.WaitGroup) error {
	return prometheus.Setup(ctx, wg, config.C.Prometheus.Bind)
}

func startSimulator(ctx context.Context, wg *sync.WaitGroup) error {
	if err := simulator.Start(ctx, wg, config.C); err != nil {
		log.Fatalf("Failed to start simulator: %v", err)
		return err
	}
	log.Info("Simulator started successfully")
	return nil
}
