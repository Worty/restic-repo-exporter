package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus"
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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	exp := resticrepoexporter.Exporter{
		RootPath: *repoPath,
	}
	go exp.Scan(ctx)

	registry := prometheus.NewRegistry()
	registry.MustRegister(&exp)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))

	log.Printf("Starting exporter on %s\n", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
