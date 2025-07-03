package resticrepoexporter

import (
	"context"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

type Exporter struct {
	repoPath string
	repos    sync.Map // map[string]*Repo
}

func NewExporter(ctx context.Context, path string) *Exporter {
	exp := Exporter{
		repoPath: path,
	}

	go exp.Scan(ctx)

	return &exp
}

func (e *Exporter) Scan(ctx context.Context) error {
	for {
		dirs, err := os.ReadDir(e.repoPath)
		if err != nil {
			log.Printf("Error reading repo path %s: %v", e.repoPath, err)
		}
		for _, dir := range dirs {
			if err := ctx.Err(); err != nil {
				continue
			}
			if !dir.IsDir() {
				continue
			}
			name := dir.Name()
			if _, ok := e.repos.Load(name); ok {
				// Repo already exists, skip it
				continue
			}
			pw := os.Getenv("RESTIC_PASSWORD")
			if ov := os.Getenv("RESTIC_PASSWORD_" + name); ov != "" {
				pw = ov
			}
			repo := &Repo{
				Name:       name,
				Repository: path.Join(e.repoPath, name),
				Password:   pw,
			}
			if _, err := repo.Check(); err != nil {
				log.Printf("Error checking repo %s: %v", name, err)
				continue
			}
			log.Printf("Found new repo: %s", name)
			go repo.Scrape(ctx)
			e.repos.Store(name, repo)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(20 * time.Second):

		}
	}
}
