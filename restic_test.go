package resticrepoexporter

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var testLoc *time.Location = time.FixedZone("CEST", 2*60*60)

func TestRepoGroupedSnapshots(t *testing.T) {
	repo := Repo{
		Name:       "test-repo",
		Repository: "./testdata/repo-with-tags",
		Password:   "abc123",
	}

	want := CheckResult{
		MessageType:        "summary",
		NumErrors:          0,
		BrokenPacks:        nil,
		SuggestRepairIndex: false,
		SuggestPrune:       false,
	}

	got, err := repo.Check()
	if err != nil {
		t.Fatalf("Expected check status to be not locked, err = %v", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("CheckResult mismatch (-got +want):\n%s", diff)
	}

	gotSnapshots, err := repo.LatestSnapshots("host,tags")
	if err != nil {
		t.Fatalf("Expected to get latest snapshots, got error: %v", err)
	}

	wantSnapshots := []GroupedSnapshot{
		{
			GroupKey: GroupKey{
				Hostname: "test-host",
				Paths:    nil,
				Tags:     []string{"test-tag"},
			},
			Snapshots: []Snapshot{
				{
					Time:           time.Date(2025, 7, 2, 19, 46, 45, 543749354, testLoc),
					Tree:           "5a6aabcc299676636b4ef62302d6b4f5363800a6f7f305b5336cd39b6697fd9a",
					Paths:          []string{"/home/worty/Projekte/restic-repo-exporter/go.mod"},
					Hostname:       "test-host",
					Username:       "worty",
					UID:            1000,
					GID:            1000,
					Tags:           []string{"test-tag"},
					ProgramVersion: "restic 0.18.0",
					Summary: Summary{
						BackupStart:         time.Date(2025, 7, 2, 19, 46, 45, 543749354, testLoc),
						BackupEnd:           time.Date(2025, 7, 2, 19, 46, 46, 266989420, testLoc),
						FilesNew:            1,
						FilesChanged:        0,
						FilesUnmodified:     0,
						DirsNew:             0,
						DirsChanged:         0,
						DirsUnmodified:      0,
						DataBlobs:           1,
						TreeBlobs:           1,
						DataAdded:           435,
						DataAddedPacked:     440,
						TotalFilesProcessed: 1,
						TotalBytesProcessed: 56,
					},
					ID:      "49053231c3f32332fed6144226222c724693def1bdbf34a60dab00e0715ad3c5",
					ShortID: "49053231",
				},
			},
		},
		{
			GroupKey: GroupKey{
				Hostname: "test-host",
				Paths:    nil,
				Tags:     []string{"test-tag-latest"},
			},
			Snapshots: []Snapshot{
				{
					Time:           time.Date(2025, 7, 2, 19, 46, 53, 759177012, testLoc),
					Parent:         "49053231c3f32332fed6144226222c724693def1bdbf34a60dab00e0715ad3c5",
					Tree:           "5a6aabcc299676636b4ef62302d6b4f5363800a6f7f305b5336cd39b6697fd9a",
					Paths:          []string{"/home/worty/Projekte/restic-repo-exporter/go.mod"},
					Hostname:       "test-host",
					Username:       "worty",
					UID:            1000,
					GID:            1000,
					Tags:           []string{"test-tag-latest"},
					ProgramVersion: "restic 0.18.0",
					Summary: Summary{
						BackupStart:         time.Date(2025, 7, 2, 19, 46, 53, 759177012, testLoc),
						BackupEnd:           time.Date(2025, 7, 2, 19, 46, 54, 399548120, testLoc),
						FilesNew:            0,
						FilesChanged:        0,
						FilesUnmodified:     1,
						DirsNew:             0,
						DirsChanged:         0,
						DirsUnmodified:      0,
						DataBlobs:           0,
						TreeBlobs:           0,
						DataAdded:           0,
						DataAddedPacked:     0,
						TotalFilesProcessed: 1,
						TotalBytesProcessed: 56,
					},
					ID:      "d79f2d3d280417dbbc9b5c44467592a8de11331d81136f505b58e235eb2c6161",
					ShortID: "d79f2d3d",
				},
			},
		},
	}

	SortGroupedSnapshots(gotSnapshots)
	SortGroupedSnapshots(wantSnapshots)

	if diff := cmp.Diff(gotSnapshots, wantSnapshots); diff != "" {
		t.Errorf("Snapshots mismatch (-got +want):\n%s", diff)
	}

	wantStats := rawDataStats{
		TotalSize:              358,
		TotalUncompressedSize:  499,
		CompressionRatio:       1.393854748603352,
		CompressionProgress:    100,
		CompressionSpaceSaving: 28.256513026052097,
		TotalBlobCount:         2,
		SnapshotsCount:         2,
	}

	gotStats, err := repo.RawStats([]string{})
	if err != nil {
		t.Fatalf("Expected to get raw data stats, got error: %v", err)
	}

	if diff := cmp.Diff(gotStats, wantStats); diff != "" {
		t.Errorf("Raw data stats mismatch (-got +want):\n%s", diff)
	}

	wantStats = rawDataStats{
		TotalSize:              358,
		TotalUncompressedSize:  499,
		CompressionRatio:       1.393854748603352,
		CompressionProgress:    100,
		CompressionSpaceSaving: 28.256513026052097,
		TotalBlobCount:         2,
		SnapshotsCount:         1,
	}

	gotStats, err = repo.RawStats([]string{"test-tag"})
	if err != nil {
		t.Fatalf("Expected to get raw data stats with tag, got error: %v", err)
	}

	if diff := cmp.Diff(gotStats, wantStats); diff != "" {
		t.Errorf("Raw data stats with tag mismatch (-got +want):\n%s", diff)
	}
}

func TestRepoLocked(t *testing.T) {
	repo := Repo{
		Name:       "test-repo-locked",
		Repository: "./testdata/repo-locked",
		Password:   "abc123",
	}

	_, err := repo.Check()
	if err == nil {
		t.Fatalf("Expected check status to be locked, got no error")
	}
}
