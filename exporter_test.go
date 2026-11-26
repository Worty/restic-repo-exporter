package resticrepoexporter

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMetrics(t *testing.T) {
	exporter := Exporter{
		Path:                  "testdata",
		ScrapeIntervalSeconds: 1,
	}

	os.Setenv("RESTIC_PASSWORD_locked-repo", "abc123")
	os.Setenv("RESTIC_PASSWORD_repo-with-tags", "abc123")
	os.Setenv("RESTIC_PASSWORD_repo-with-broken-data", "abc123")

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	exporter.Scan(ctx)

	got := testutil.ToFloat64(numSnapshots.WithLabelValues("locked-repo", "Wortys-Thinkpad", ""))
	want := 2.0

	if want != got {
		t.Errorf("unexpected number of snapshots: got %f, want %f", got, want)
	}
}

func TestMetricsStripTags(t *testing.T) {
	exporter := Exporter{
		Path:                  "testdata",
		ScrapeIntervalSeconds: 1,
		StripTags:             []string{"dont-prune"},
	}

	os.Setenv("RESTIC_PASSWORD_locked-repo", "abc123")
	os.Setenv("RESTIC_PASSWORD_repo-with-tags", "abc123")
	os.Setenv("RESTIC_PASSWORD_repo-with-broken-data", "abc123")

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	exporter.Scan(ctx)

	// The snapshot tagged "dont-prune,test-tag" must still be reported, but with
	// the stripped tag removed from the label: tag="test-tag".
	if got, ok := gaugeValue(t, "restic_repo_number_of_snapshots", prometheus.Labels{
		"repo":     "repo-with-tags",
		"hostname": "test-host",
		"tag":      "test-tag",
	}); !ok || got != 1.0 {
		t.Errorf("unexpected number of snapshots for tag test-tag: got %v (found: %v), want 1", got, ok)
	}

	// The stripped tag itself must not be reported as a label.
	if _, ok := gaugeValue(t, "restic_repo_number_of_snapshots", prometheus.Labels{
		"repo":     "repo-with-tags",
		"hostname": "test-host",
		"tag":      "dont-prune_test-tag",
	}); ok {
		t.Error("snapshot still reported with the stripped tag label dont-prune_test-tag")
	}
}

// gaugeValue returns the value of the gauge metric with the given name and
// label values, or false if no such metric exists.
func gaugeValue(t *testing.T, name string, labels prometheus.Labels) (float64, bool) {
	t.Helper()

	mfs, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}

	for _, mf := range mfs {
		if mf.GetName() != name {
			continue
		}
		for _, m := range mf.GetMetric() {
			if len(m.GetLabel()) != len(labels) {
				continue
			}
			matches := true
			for _, label := range m.GetLabel() {
				if labels[label.GetName()] != label.GetValue() {
					matches = false
					break
				}
			}
			if matches {
				return m.GetGauge().GetValue(), true
			}
		}
	}

	return 0, false
}
