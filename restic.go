package resticrepoexporter

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	scrapeErr = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "restic",
			Subsystem: "repo",
			Name:      "scrape_errors_total",
			Help:      "Total number of errors encountered while scraping restic repository",
		},
		[]string{"repo", "action"},
	)
	scrapeDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "restic",
			Subsystem: "repo",
			Name:      "scrape_duration_seconds",
			Help:      "Duration of the last scrape of the restic repository",
			Buckets:   prometheus.ExponentialBucketsRange(0.1, 15, 10),
		},
		[]string{"repo"},
	)
	numSnapshotsDesc = prometheus.NewDesc(
		prometheus.BuildFQName("restic", "repo", "number_of_snapshots"),
		"Total number of snapshots in the repository by hostname and tag",
		[]string{"repo", "hostname", "tag"},
		nil,
	)
	lastSnapshotTimestampDesc = prometheus.NewDesc(
		prometheus.BuildFQName("restic", "repo", "last_snapshot_timestamp"),
		"Timestamp of the last snapshot in the repository by hostname and tag",
		[]string{"repo", "hostname", "tag"},
		nil,
	)
	numRepoErrorsDesc = prometheus.NewDesc(
		prometheus.BuildFQName("restic", "repo", "num_errors"),
		"Total number of errors found in the repository during check",
		[]string{"repo"},
		nil,
	)
	suggestRepairIndexDesc = prometheus.NewDesc(
		prometheus.BuildFQName("restic", "repo", "suggest_repair_index"),
		"Whether the repository suggests repairing the index",
		[]string{"repo"},
		nil,
	)
	suggestPruneIndexDesc = prometheus.NewDesc(
		prometheus.BuildFQName("restic", "repo", "suggest_prune"),
		"Whether the repository suggests pruning",
		[]string{"repo"},
		nil,
	)
	totalRepoSizeDesc = prometheus.NewDesc(
		prometheus.BuildFQName("restic", "repo", "total_size_bytes"),
		"Total size of the repository in bytes",
		[]string{"repo"},
		nil,
	)
	totalUncompressedSizeDesc = prometheus.NewDesc(
		prometheus.BuildFQName("restic", "repo", "total_uncompressed_size_bytes"),
		"Total uncompressed size of the repository in bytes",
		[]string{"repo"},
		nil,
	)
	compressionRatioDesc = prometheus.NewDesc(
		prometheus.BuildFQName("restic", "repo", "compression_ratio"),
		"Compression ratio of the repository",
		[]string{"repo"},
		nil,
	)
	compressionProgressDesc = prometheus.NewDesc(
		prometheus.BuildFQName("restic", "repo", "compression_progress_percent"),
		"Compression progress of the repository in percent",
		[]string{"repo"},
		nil,
	)
	compressionSpaceSavingDesc = prometheus.NewDesc(
		prometheus.BuildFQName("restic", "repo", "compression_space_saving_percent"),
		"Compression space saving of the repository in percent",
		[]string{"repo"},
		nil,
	)
	totalBlobCountDesc = prometheus.NewDesc(
		prometheus.BuildFQName("restic", "repo", "total_blob_count"),
		"Total number of blobs in the repository",
		[]string{"repo"},
		nil,
	)
	totalSnapshotsCountDesc = prometheus.NewDesc(
		prometheus.BuildFQName("restic", "repo", "total_snapshots_count"),
		"Total number of snapshots in the repository",
		[]string{"repo"},
		nil,
	)
)

type Repo struct {
	Name       string
	Repository string
	Password   string
}

func (r *Repo) Describe(ch chan<- *prometheus.Desc) {
	ch <- numSnapshotsDesc
	ch <- numRepoErrorsDesc
	ch <- suggestRepairIndexDesc
	ch <- suggestPruneIndexDesc
}

func (r *Repo) Collect(ch chan<- prometheus.Metric) {
	timeStart := time.Now()
	defer func() {
		scrapeDuration.WithLabelValues(r.Name).Observe(time.Since(timeStart).Seconds())
	}()

	if check, err := r.Check(); err == nil {
		ch <- prometheus.MustNewConstMetric(numRepoErrorsDesc, prometheus.GaugeValue, float64(check.NumErrors), r.Name)
		ch <- prometheus.MustNewConstMetric(suggestRepairIndexDesc, prometheus.GaugeValue, boolToFloat(check.SuggestRepairIndex), r.Name)
		ch <- prometheus.MustNewConstMetric(suggestPruneIndexDesc, prometheus.GaugeValue, boolToFloat(check.SuggestPrune), r.Name)
	} else {
		scrapeErr.WithLabelValues(r.Name, "check").Inc()
		return
	}

	rawStats, err := r.RawStats(nil)
	if err != nil {
		scrapeErr.WithLabelValues(r.Name, "raw-stats").Inc()
		return
	}
	ch <- prometheus.MustNewConstMetric(totalRepoSizeDesc, prometheus.GaugeValue, float64(rawStats.TotalSize), r.Name)
	ch <- prometheus.MustNewConstMetric(totalUncompressedSizeDesc, prometheus.GaugeValue, float64(rawStats.TotalUncompressedSize), r.Name)
	ch <- prometheus.MustNewConstMetric(compressionRatioDesc, prometheus.GaugeValue, rawStats.CompressionRatio, r.Name)
	ch <- prometheus.MustNewConstMetric(compressionProgressDesc, prometheus.GaugeValue, float64(rawStats.CompressionProgress)/100.0, r.Name)
	ch <- prometheus.MustNewConstMetric(compressionSpaceSavingDesc, prometheus.GaugeValue, (rawStats.CompressionSpaceSaving)/100.0, r.Name)
	ch <- prometheus.MustNewConstMetric(totalBlobCountDesc, prometheus.GaugeValue, float64(rawStats.TotalBlobCount), r.Name)
	ch <- prometheus.MustNewConstMetric(totalSnapshotsCountDesc, prometheus.GaugeValue, float64(rawStats.SnapshotsCount), r.Name)

	snapshots, err := r.LatestSnapshots("host,tags")
	if err != nil {
		scrapeErr.WithLabelValues(r.Name, "latest-snapshots").Inc()
		return
	}
	for _, group := range snapshots {
		for _, snapshot := range group.Snapshots {
			ch <- prometheus.MustNewConstMetric(numSnapshotsDesc, prometheus.GaugeValue, float64(len(group.Snapshots)), r.Name, group.GroupKey.Hostname, strings.Join(group.GroupKey.Tags, "_"))
			ch <- prometheus.MustNewConstMetric(lastSnapshotTimestampDesc, prometheus.GaugeValue, float64(snapshot.Time.Unix()), r.Name, group.GroupKey.Hostname, strings.Join(group.GroupKey.Tags, "_"))
		}
	}
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
	cmd := exec.Command("restic", "--quiet", "--no-lock", "--json")
	cmd.Args = append(cmd.Args, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "RESTIC_REPOSITORY="+r.Repository)
	cmd.Env = append(cmd.Env, "RESTIC_PASSWORD="+r.Password)

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

func (r *Repo) LatestSnapshots(groupBy string) (gr []GroupedSnapshot, err error) {
	output, err := r.exec("snapshots", "--latest=1", "--group-by="+groupBy)
	if err != nil {
		return []GroupedSnapshot{}, err
	}

	return gr, json.Unmarshal(output, &gr)
}

type CheckResult struct {
	MessageType        string   `json:"message_type"`
	NumErrors          int      `json:"num_errors"`
	BrokenPacks        []string `json:"broken_packs"`
	SuggestRepairIndex bool     `json:"suggest_repair_index"`
	SuggestPrune       bool     `json:"suggest_prune"`
}

func (r *Repo) Check() (cr CheckResult, err error) {
	o, err := r.exec("check")
	if err != nil {
		return CheckResult{}, err
	}

	return cr, json.Unmarshal(o, &cr)
}

func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}
