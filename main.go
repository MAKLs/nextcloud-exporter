package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MAKLs/nextcloud-exporter/client"
	"github.com/MAKLs/nextcloud-exporter/config"
	"github.com/MAKLs/nextcloud-exporter/exporter"
	"github.com/MAKLs/nextcloud-exporter/metrics"
	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	shutdownTimeout = 5 * time.Second
)

var (
	ncRegistry = metrics.ExporterRegistry
	ncExporter *exporter.NCExporter
	mux        *http.ServeMux
)

func healthz() http.Handler {
	health := struct {
		Status string `json:"status"`
	}{Status: "UP"}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := json.Marshal(health)
		if err != nil {
			log.Fatalf("failed to serialize health status: %s", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
}

func stop(serverChan <-chan *http.Server, errorChan chan<- error) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer func() {
		cancel()
	}()

	server := <-serverChan
	log.Println("stopping server")
	errorChan <- server.Shutdown(ctx)
	ncRegistry.Unregister(ncExporter)
}

func start(serverChan chan<- *http.Server, errorChan chan<- error) {
	appConfig := config.GetConfig()
	ncClient := client.NewNCClient(&appConfig.URL, appConfig.Token)
	ncExporter = exporter.NewNCExporter(ncClient, appConfig.ExcludePHP, appConfig.ExcludeStrings, appConfig.FilterMetrics)
	ncRegistry.MustRegister(ncExporter)
	server := &http.Server{Handler: mux, Addr: fmt.Sprintf(":%d", appConfig.Port)}
	serverChan <- server
	log.Printf("starting server at :%d", appConfig.Port)
	errorChan <- server.ListenAndServe()
}

func restart(serverChan chan *http.Server, errorChan chan error) {
	stop(serverChan, errorChan)
	start(serverChan, errorChan)
}

func main() {
	// Prepare channels
	serverChan := make(chan *http.Server, 1)
	reloadChan := make(chan fsnotify.Event)
	// Buffer for 1 start and 1 shutdown error so shutdown and reload signals aren't blocked
	errorChan := make(chan error, 2)
	shutdownChan := make(chan os.Signal, 1)
	doneChan := make(chan bool)

	config.Notify(reloadChan)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Prepare endpoints
	mux = http.NewServeMux()
	mux.Handle("/healthz", healthz())
	mux.Handle("/metrics", promhttp.HandlerFor(metrics.ExporterRegistry, promhttp.HandlerOpts{}))

	// Initial start
	go start(serverChan, errorChan)

	go func() {
		for {
			select {
			// Watch for config changes to reload exporter
			case ev := <-reloadChan:
				log.Printf("detected %s to config \"%s\", restarting", ev.Op, ev.Name)
				go restart(serverChan, errorChan)
			// Watch for non-recoverable errors
			case err := <-errorChan:
				switch err {
				case nil:
				case http.ErrServerClosed:
				case context.DeadlineExceeded:
					log.Printf("failed to stop server within deadline (%d s)", shutdownTimeout)
				default:
					log.Printf("unexpected error: %v", err)
					close(doneChan)
				}
			// Watch for shutdown signals
			case <-shutdownChan:
				log.Println("received SIGINT")
				stop(serverChan, errorChan)
				doneChan <- true
			}
		}
	}()

	<-doneChan
}
