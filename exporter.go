package resticrepoexporter

import (
	"context"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

type Exporter struct {
	RootPath string
	Repos    sync.Map // map[string]*Repo
}

// Collect implements [prometheus.Collector]
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Only use [errgroup.Group] to limit concurrency, not to collect errors
	var eg errgroup.Group
	eg.SetLimit(runtime.NumCPU() * 2)
	e.Repos.Range(func(key, value any) bool {
		r := value.(*Repo)
		eg.Go(func() error {
			r.Collect(ch)
			return nil
		})
		return true
	})
	eg.Wait()
	scrapeDuration.Collect(ch)
	scrapeErr.Collect(ch)
}

// Describe implements [prometheus.Collector]
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.Repos.Range(func(key, value any) bool {
		// name := key.(string)
		r := value.(*Repo)
		r.Describe(ch)
		return true
	})
	scrapeErr.Describe(ch)
	scrapeDuration.Describe(ch)
}

func (e *Exporter) Scan(ctx context.Context) error {
	for {
		dirs, err := os.ReadDir(e.RootPath)
		if err != nil {
			log.Printf("Error reading repo path %s: %v", e.RootPath, err)
		}
		for _, dir := range dirs {
			if !dir.IsDir() {
				continue
			}
			name := dir.Name()
			if _, ok := e.Repos.Load(name); ok {
				// Repo already exists, skip it
				continue
			}
			pw := os.Getenv("RESTIC_PASSWORD")
			if ov := os.Getenv("RESTIC_PASSWORD_" + name); ov != "" {
				pw = ov
			}
			repo := &Repo{
				Name:       name,
				Repository: path.Join(e.RootPath, name),
				Password:   pw,
			}
			_, err := repo.Check()
			if err != nil {
				log.Printf("Error checking repo %s: %v", name, err)
				continue
			}
			e.Repos.Store(name, repo)
			log.Printf("Found new repo: %s", name)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(20 * time.Second)
		}
	}
}
