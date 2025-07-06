package resticrepoexporter

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	scrapeErr = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "restic",
			Subsystem: "repo",
			Name:      "scrape_errors_total",
			Help:      "Total number of errors encountered while scraping restic repository",
		},
		[]string{"repo", "action"},
	)
	scrapeDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "restic",
			Subsystem: "repo",
			Name:      "scrape_duration_seconds",
			Buckets:   prometheus.ExponentialBucketsRange(0.1, 60, 10),
		},
		[]string{"repo", "action"},
	)
	numSnapshots = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("restic", "repo", "number_of_snapshots"),
			Help: "Total number of snapshots in the repository by hostname and tag",
		},
		[]string{"repo", "hostname", "tag"},
	)
	lastSnapshotCreationDuration = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("restic", "repo", "last_snapshot_creation_seconds"),
			Help: "Time it took to create the last snapshot",
		},
		[]string{"repo", "hostname", "tag"},
	)
	lastSnapshotTimestamp = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("restic", "repo", "last_snapshot_timestamp"),
			Help: "Timestamp of the last snapshot in the repository by hostname and tag",
		},
		[]string{"repo", "hostname", "tag"},
	)
	numRepoErrors = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("restic", "repo", "num_errors"),
			Help: "Total number of errors found in the repository during check",
		},
		[]string{"repo"},
	)
	suggestRepairIndex = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("restic", "repo", "suggest_repair_index"),
			Help: "Whether the repository suggests repairing the index",
		},
		[]string{"repo"},
	)
	suggestPrune = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("restic", "repo", "suggest_prune"),
			Help: "Whether the repository suggests pruning",
		},
		[]string{"repo"},
	)
	totalRepoSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("restic", "repo", "total_size_bytes"),
			Help: "Total size of the repository in bytes",
		},
		[]string{"repo"},
	)
	totalUncompressedSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("restic", "repo", "total_uncompressed_size_bytes"),
			Help: "Total uncompressed size of the repository in bytes",
		},
		[]string{"repo"},
	)
	compressionRatio = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("restic", "repo", "compression_ratio"),
			Help: "Compression ratio of the repository",
		},
		[]string{"repo"},
	)
	compressionProgress = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("restic", "repo", "compression_progress_percent"),
			Help: "Compression progress of the repository in percent",
		},
		[]string{"repo"},
	)
	compressionSpaceSaving = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("restic", "repo", "compression_space_saving_percent"),
			Help: "Compression space saving of the repository in percent",
		},
		[]string{"repo"},
	)
	totalBlobCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("restic", "repo", "total_blob_count"),
			Help: "Total number of blobs in the repository",
		},
		[]string{"repo"},
	)

	totalSnapshotsCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName("restic", "repo", "total_snapshots_count"),
			Help: "Total number of snapshots in the repository",
		},
		[]string{"repo"},
	)
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
				return fs.SkipAll
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
			if ov := os.Getenv("RESTIC_PASSWORD_" + dir.Name()); ov != "" {
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
