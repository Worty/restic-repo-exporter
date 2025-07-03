package main

import (
	"context"
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
	flag.Parse()

	if *repoPath == "" {
		log.Fatal("repo-path must be specified")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, os.Kill)
	defer cancel()

	resticrepoexporter.NewExporter(ctx, *repoPath)

	// Create a new ServeMux for custom routing
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	log.Printf("Starting Prometheus metrics server on %s", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, mux); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
