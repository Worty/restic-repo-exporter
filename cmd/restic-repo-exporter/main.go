package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	resticrepoexporter "github.com/worty/restic-repo-exporter"
)

func main() {
	listenAddr := flag.String("listen-address", ":9100", "The address to listen on for HTTP requests.")
	repoPath := flag.String("repo-path", "", "Path to the directory containing restic repositories.")
	scrapeInterval := flag.Int64("scrape-interval", 30, "Base scrape interval in seconds. A random interval of the same amount will be added on top.")
	flag.Parse()

	if *repoPath == "" {
		log.Fatal("repo-path must be specified")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, os.Kill)
	defer cancel()

	resticrepoexporter.NewExporter(ctx, *repoPath, *scrapeInterval)

	// Create a new ServeMux for custom routing
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := http.Server{Addr: *listenAddr, Handler: mux}

	context.AfterFunc(ctx, func() {
		if err := srv.Close(); err != nil {
			log.Printf("srv.Close() err = %v", err)
		}
	})

	log.Printf("Starting Prometheus metrics server on %s", srv.Addr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
