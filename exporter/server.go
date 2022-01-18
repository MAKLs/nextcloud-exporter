package exporter

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MAKLs/nextcloud-exporter/client"
	"github.com/MAKLs/nextcloud-exporter/config"
)

const (
	shutdownTimeout = 5 * time.Second
)

// Stop reads an http server from the server channel, stops it and reports the error to the error channel
//
// The server channel is written to in the Start method after starting the http server.
func Stop(serverChan <-chan *http.Server, errorChan chan<- error) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	server := <-serverChan
	log.Println("stopping server")
	errorChan <- server.Shutdown(ctx)
}

// Start configures the exporter with the current configuration values, starts an http server with the provided mux,
// writes the server to the server channel and reports errors to the error channel.
//
// Stop reads from the server channel to stop the server later.
func Start(serveMux *http.ServeMux, serverChan chan<- *http.Server, errorChan chan<- error) {
	appConfig := config.GetConfig()
	ncClient := client.NewNCClient(&appConfig.URL, appConfig.Token)
	ConfigureExporter(ncClient, appConfig.ExcludePHP, appConfig.FilterMetrics)
	server := &http.Server{Handler: serveMux, Addr: fmt.Sprintf(":%d", appConfig.Port)}
	serverChan <- server
	log.Printf("starting server at :%d", appConfig.Port)
	errorChan <- server.ListenAndServe()
}

// Restart is a shorthand wrapper around Stop and Start.
func Restart(serveMux *http.ServeMux, serverChan chan *http.Server, errorChan chan error) {
	Stop(serverChan, errorChan)
	Start(serveMux, serverChan, errorChan)
}
