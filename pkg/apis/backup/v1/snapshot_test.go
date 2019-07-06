package v1

import (
	"testing"

	backup "github.com/linuxuser586/apis/grpc/cassandra/backup/v1"
	"github.com/spf13/afero"
)

func TestSnapshot(t *testing.T) {
	fs = afero.NewMemMapFs()
	loadFiles(t)
	ds := &backup.Downstream{
		Type: backup.Downstream_SNAPSHOT,
	}
	s := backupService{}
	service := newFakeService()
	if err := s.StreamFrom(ds, service); err != nil {
		t.Error(err)
	}
	close(service.in)
	dir := afero.NewBasePathFs(fs, tmpDir)
	if err := recv(t, dir, service); err != nil {
		t.Fatal(err)
	}
	td := []struct {
		path string
		file string
		hash string
	}{
		{"keyspace1/table_1_12345/snapshots/2544512/", "data-1.txt", "1YUjpA=="},
	}
	for _, tt := range td {
		meta := &backup.Upstream_Meta{
			FileName: tt.file,
			Hash:     tt.hash,
			Path:     tt.path,
		}
		verifyFile(t, dir, meta)
	}
}
