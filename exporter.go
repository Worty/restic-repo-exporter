package resticrepoexporter

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

type Exporter struct {
	repoPath              string
	scrapeIntervalSeconds int64
	repos                 sync.Map // map[string]*Repo
}

func NewExporter(ctx context.Context, path string, scrapeIntervalSeconds int64) *Exporter {
	exp := Exporter{
		repoPath:              path,
		scrapeIntervalSeconds: scrapeIntervalSeconds,
	}

	go exp.Scan(ctx)

	return &exp
}

func (e *Exporter) Scan(ctx context.Context) error {
	for {
		fs.WalkDir(os.DirFS(e.repoPath), ".", func(dirPath string, dir fs.DirEntry, err error) error {
			if err := ctx.Err(); err != nil {
				return nil
			}
			if !dir.IsDir() {
				return nil
			}

			if _, ok := e.repos.Load(dirPath); ok {
				return fs.SkipDir
			}

			if _, err := os.Stat(path.Join(e.repoPath, dirPath, "config")); os.IsNotExist(err) {
				return nil
			}

			pw := os.Getenv("RESTIC_PASSWORD")
			if ov := os.Getenv("RESTIC_PASSWORD_" + dirPath); ov != "" {
				pw = ov
			}
			repo := &Repo{
				Name:       dir.Name(),
				Repository: path.Join(e.repoPath, dirPath),
				Password:   pw,
			}
			log.Printf("Found new repo: %s", dir.Name())
			if _, err := repo.Check(); err != nil {
				log.Printf("Error checking repo %s: %v", dirPath, err)
				return fs.SkipDir
			}
			go repo.Scrape(ctx, e.scrapeIntervalSeconds)
			e.repos.Store(dirPath, repo)
			return fs.SkipDir
		})

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(20 * time.Second):

		}
	}
}
