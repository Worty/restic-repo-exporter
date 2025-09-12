package resticrepoexporter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Repo struct {
	Name       string
	Repository string
	Password   string

	modTimes map[string]time.Time
}

func (r *Repo) Scrape(ctx context.Context, scrapeIntervalSeconds int64, semaphore chan struct{}) {
	for {
		// To always sleep even if we got an error
		func() {
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if !r.changed() {
				return
			}
			log.Printf("Start Scraping Repo %s", r.Name)
			timer := prometheus.NewTimer(scrapeDuration.WithLabelValues(r.Name, "config"))
			if config, err := r.Config(); err == nil {
				repoVersion.WithLabelValues(r.Name).Set(float64(config.Version))
			} else {
				scrapeErr.WithLabelValues(r.Name, "config").Inc()
			}
			timer.ObserveDuration()

			timer = prometheus.NewTimer(scrapeDuration.WithLabelValues(r.Name, "check"))
			if check, err := r.Check(); err == nil {
				numRepoErrors.WithLabelValues(r.Name).Set(float64(check.NumErrors))
				suggestPrune.WithLabelValues(r.Name).Set(boolToFloat(check.SuggestPrune))
				suggestRepairIndex.WithLabelValues(r.Name).Set(boolToFloat(check.SuggestRepairIndex))
			} else {
				scrapeErr.WithLabelValues(r.Name, "check").Inc()
				return
			}
			timer.ObserveDuration()

			timer = prometheus.NewTimer(scrapeDuration.WithLabelValues(r.Name, "raw-stats"))
			rawStats, err := r.RawStats(nil)
			if err != nil {
				scrapeErr.WithLabelValues(r.Name, "raw-stats").Inc()
				return
			}
			totalRepoSize.WithLabelValues(r.Name).Set(float64(rawStats.TotalSize))
			totalUncompressedSize.WithLabelValues(r.Name).Set(float64(rawStats.TotalUncompressedSize))
			compressionRatio.WithLabelValues(r.Name).Set(rawStats.CompressionRatio)
			compressionProgress.WithLabelValues(r.Name).Set(float64(rawStats.CompressionProgress) / 100.0)
			compressionSpaceSaving.WithLabelValues(r.Name).Set(float64(rawStats.CompressionSpaceSaving) / 100.0)
			totalBlobCount.WithLabelValues(r.Name).Set(float64(rawStats.TotalBlobCount))
			totalSnapshotsCount.WithLabelValues(r.Name).Set(float64(rawStats.SnapshotsCount))
			timer.ObserveDuration()

			timer = prometheus.NewTimer(scrapeDuration.WithLabelValues(r.Name, "snapshots"))
			groups, err := r.Snapshots("host,tags")
			if err != nil {
				scrapeErr.WithLabelValues(r.Name, "snapshots").Inc()
				return
			}
			for _, group := range groups {
				tags := strings.Join(group.GroupKey.Tags, "_")
				numSnapshots.WithLabelValues(r.Name, group.GroupKey.Hostname, tags).Set(float64(len(group.Snapshots)))
				if len(group.Snapshots) > 0 {
					lastSnapshot := group.Snapshots[0]
					lastSnapshotTimestamp.WithLabelValues(r.Name, group.GroupKey.Hostname, tags).Set(float64(lastSnapshot.Time.Unix()))
					lastSnapshotCreationDuration.WithLabelValues(r.Name, group.GroupKey.Hostname, tags).Set((lastSnapshot.Summary.BackupEnd.Sub(lastSnapshot.Summary.BackupStart)).Seconds())
				}
			}
			timer.ObserveDuration()
		}()

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(rand.Int64N(scrapeIntervalSeconds)+scrapeIntervalSeconds) * time.Second):

		}
	}
}

func (r *Repo) changed() (changed bool) {
	if r.modTimes == nil {
		r.modTimes = make(map[string]time.Time, 7)
	}

	dir, err := os.ReadDir(r.Repository)
	if err != nil {
		log.Printf("error checking if repo chaged: os.ReadDir(%q) = %v", r.Repository, err)
		return true
	}

	for _, d := range dir {
		i, err := d.Info()
		if err != nil {
			log.Printf("error checking if repo chaged: %s.Info() = %v", d.Name(), err)
			return true
		}

		if r.modTimes[i.Name()].Before(i.ModTime()) {
			r.modTimes[i.Name()] = i.ModTime()
			changed = true
		}
	}
	return changed
}

type rawDataStats struct {
	TotalSize              int     `json:"total_size"`
	TotalUncompressedSize  int     `json:"total_uncompressed_size"`
	CompressionRatio       float64 `json:"compression_ratio"`
	CompressionProgress    int     `json:"compression_progress"`
	CompressionSpaceSaving float64 `json:"compression_space_saving"`
	TotalBlobCount         int     `json:"total_blob_count"`
	SnapshotsCount         int     `json:"snapshots_count"`
}

func (r *Repo) RawStats(tags []string) (stats rawDataStats, err error) {
	args := []string{"stats", "--mode", "raw-data"}
	if len(tags) > 0 {
		args = append(args, "--tag="+strings.Join(tags, ","))
	}
	o, err := r.exec(args...)
	if err != nil {
		return rawDataStats{}, fmt.Errorf("error executing stats command: %w", err)
	}
	return stats, json.Unmarshal(o, &stats)
}

func (r *Repo) exec(args ...string) ([]byte, error) {
	cmd := exec.Command("restic", "-r", r.Repository, "--quiet", "--no-lock", "--json")
	cmd.Args = append(cmd.Args, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "RESTIC_PASSWORD="+r.Password)
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOMAXPROCS=%d", runtime.NumCPU()/4)) // One Restic Procress should only take 25% of the CPU

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error executing command %s: %w output: %s", cmd, err, out)
	}

	return out, nil
}

type GroupedSnapshot struct {
	GroupKey  GroupKey   `json:"group_key"`
	Snapshots []Snapshot `json:"snapshots"`
}

func SortGroupedSnapshots(groups []GroupedSnapshot) {
	for i := range groups {
		slices.SortStableFunc(groups[i].Snapshots, func(a, b Snapshot) int {
			if a.Time.After(b.Time) {
				return -1
			}
			if a.Time.Before(b.Time) {
				return 1
			}
			return 0
		})
	}

	slices.SortStableFunc(groups, func(a, b GroupedSnapshot) int {
		var aTime, bTime time.Time
		if len(a.Snapshots) > 0 {
			aTime = a.Snapshots[0].Time
		}
		if len(b.Snapshots) > 0 {
			bTime = b.Snapshots[0].Time
		}
		if aTime.After(bTime) {
			return -1
		}
		if aTime.Before(bTime) {
			return 1
		}
		return 0
	})
}

type GroupKey struct {
	Hostname string   `json:"hostname"`
	Paths    []string `json:"paths"`
	Tags     []string `json:"tags"`
}

type Snapshot struct {
	Time           time.Time `json:"time"`
	Tree           string    `json:"tree"`
	Paths          []string  `json:"paths"`
	Hostname       string    `json:"hostname"`
	Username       string    `json:"username"`
	UID            int       `json:"uid"`
	GID            int       `json:"gid"`
	Tags           []string  `json:"tags"`
	ProgramVersion string    `json:"program_version"`
	Summary        Summary   `json:"summary"`
	ID             string    `json:"id"`
	ShortID        string    `json:"short_id"`
	Parent         string    `json:"parent,omitempty"`
}

type Summary struct {
	BackupStart         time.Time `json:"backup_start"`
	BackupEnd           time.Time `json:"backup_end"`
	FilesNew            int       `json:"files_new"`
	FilesChanged        int       `json:"files_changed"`
	FilesUnmodified     int       `json:"files_unmodified"`
	DirsNew             int       `json:"dirs_new"`
	DirsChanged         int       `json:"dirs_changed"`
	DirsUnmodified      int       `json:"dirs_unmodified"`
	DataBlobs           int       `json:"data_blobs"`
	TreeBlobs           int       `json:"tree_blobs"`
	DataAdded           int       `json:"data_added"`
	DataAddedPacked     int       `json:"data_added_packed"`
	TotalFilesProcessed int       `json:"total_files_processed"`
	TotalBytesProcessed int       `json:"total_bytes_processed"`
}

func (r *Repo) Snapshots(groupBy string) (gr []GroupedSnapshot, err error) {
	output, err := r.exec("snapshots", "--group-by="+groupBy)
	if err != nil {
		return []GroupedSnapshot{}, err
	}

	if err := json.Unmarshal(output, &gr); err != nil {
		return gr, err
	}

	SortGroupedSnapshots(gr)

	return gr, nil
}

type CheckResult struct {
	MessageType        string   `json:"message_type"`
	NumErrors          int      `json:"num_errors"`
	BrokenPacks        []string `json:"broken_packs"`
	SuggestRepairIndex bool     `json:"suggest_repair_index"`
	SuggestPrune       bool     `json:"suggest_prune"`
}

func (r *Repo) Check() (cr CheckResult, err error) {
	// From https://restic.readthedocs.io/en/latest/045_working_with_repos.html : "To reuse the existing cache, you can use the --with-cache flag"
	// Can be disabled with this env flag
	args := []string{"check"}
	if os.Getenv("CHECK_WITHOUT_CACHE") == "" {
		args = append(args, "--with-cache")
	}

	o, err := r.exec(args...)
	if err != nil {
		return CheckResult{}, err
	}

	return cr, json.Unmarshal(o, &cr)
}

type ConfigResult struct {
	Version           int    `json:"version"`
	Id                string `json:"id"`
	ChunkerPolynomial string `json:"chunker_polynomial"`
}

func (r *Repo) Config() (cr ConfigResult, err error) {
	o, err := r.exec("cat", "config")
	if err != nil {
		return ConfigResult{}, err
	}

	return cr, json.Unmarshal(o, &cr)
}

func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}
