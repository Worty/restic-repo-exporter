package resticrepoexporter

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMetrics(t *testing.T) {
	exporter := Exporter{
		repoPath:              "testdata",
		scrapeIntervalSeconds: 1,
		snapshotGroups:        "host,tags",
	}

	os.Setenv("RESTIC_PASSWORD_locked-repo", "abc123")
	os.Setenv("RESTIC_PASSWORD_repo-with-tags", "abc123")

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	exporter.Scan(ctx)

	got := testutil.ToFloat64(numSnapshots.WithLabelValues("locked-repo", "Wortys-Thinkpad", ""))
	want := 2.0

	if want != got {
		t.Errorf("unexpected number of snapshots: got %f, want %f", got, want)
	}
}
