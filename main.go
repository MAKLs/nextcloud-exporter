package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MAKLs/nextcloud-exporter/config"
	"github.com/MAKLs/nextcloud-exporter/controllers"
	"github.com/MAKLs/nextcloud-exporter/exporter"
	"github.com/fsnotify/fsnotify"
)

func main() {
	// Prepare channels
	serverChan := make(chan *http.Server, 1)
	reloadChan := make(chan fsnotify.Event)
	// Starting and stopping server ALWAYS returns an error, so
	// buffer for 1 start and 1 shutdown error; otherwise, reading will be blocked
	errorChan := make(chan error, 2)
	shutdownChan := make(chan os.Signal, 1)
	doneChan := make(chan bool)

	config.Notify(reloadChan)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Prepare endpoints
	mux := http.NewServeMux()
	mux.Handle("/healthz", controllers.HealthController())
	mux.Handle("/metrics", controllers.MetricsController())

	// Initial start
	go exporter.Start(mux, serverChan, errorChan)

	go func() {
		for {
			select {
			// Watch for config changes to reload exporter
			case ev := <-reloadChan:
				log.Printf("detected %s to config \"%s\", restarting", ev.Op, ev.Name)
				go exporter.Restart(mux, serverChan, errorChan)
			// Watch for unrecoverable errors
			case err := <-errorChan:
				switch err {
				case nil:
				case http.ErrServerClosed:
				case context.DeadlineExceeded:
					log.Printf("failed to stop server")
				default:
					log.Printf("unexpected error: %v", err)
					close(doneChan)
				}
			// Watch for shutdown signals
			case <-shutdownChan:
				log.Println("received SIGINT")
				exporter.Stop(serverChan, errorChan)
				doneChan <- true
			}
		}
	}()

	<-doneChan
}
