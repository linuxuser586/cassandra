package v1

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	backup "github.com/linuxuser586/apis/grpc/cassandra/backup/v1"
	"github.com/spf13/afero"
)

const dataDir = "testdata/data/"

func TestIncremental(t *testing.T) {
	fs = afero.NewMemMapFs()
	loadFiles(t)
	ds := &backup.Downstream{
		Type: backup.Downstream_INCREMENTAL,
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
		{"keyspace1/table_1_12345/backups/", "data-1.txt", "9wltgA=="},
		{"keyspace1/table_1_12345/backups/", "data-2.txt", "vzrddA=="},
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

func loadFiles(t *testing.T) {
	err := filepath.Walk(dataDir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			t.Fatal(err)
		}
		if !fi.IsDir() {
			wpath := dataPath + strings.TrimPrefix(path, "testdata/data/")
			b, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if err := afero.WriteFile(fs, wpath, b, 644); err != nil {
				t.Fatal(err)
			}
			t.Logf("Added: %s", wpath)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteIncremental(t *testing.T) {
	fs = afero.NewMemMapFs()
	loadFiles(t)
	service := newFakeDeleteService()
	waitc := make(chan struct{})
	go func() {
		i := 0
		for {
			resp := <-service.in
			if resp == nil {
				break
			}
			if i == 0 && resp.GetFail() {
				t.Errorf("want: fail false, got: fail %v with message: %s", resp.GetFail(), resp.GetMessage())
			} else if i > 0 && !resp.GetFail() {
				t.Error("want: fail true, got: fail false")
			}
			i++
		}
		close(waitc)
	}()
	go func() {
		service.recv <- &backup.DeleteRequest{File: "keyspace1/table_1_12345/backups/data-1.txt"}
		service.recv <- &backup.DeleteRequest{File: "keyspace1/table_1_12345/backups/data-does-not-exist.txt"}
		service.recv <- &backup.DeleteRequest{File: "keyspace1/table_1_12345/data-1.txt"}
		close(service.recv)
	}()
	s := backupService{}
	if err := s.DeleteIncremental(service); err != nil {
		t.Error(err)
	}
	close(service.in)
	<-waitc
}

func newFakeDeleteService() *fakeDeleteService {
	s := &fakeDeleteService{}
	s.in = make(chan *backup.DeleteResponse)
	s.recv = make(chan *backup.DeleteRequest)
	return s
}

type fakeDeleteService struct {
	in   chan *backup.DeleteResponse
	recv chan *backup.DeleteRequest
	fakeServerStream
}

func (f *fakeDeleteService) Send(in *backup.DeleteResponse) error {
	f.in <- in
	return nil
}

func (f *fakeDeleteService) Recv() (*backup.DeleteRequest, error) {
	v := <-f.recv
	if v == nil {
		return nil, io.EOF
	}
	return v, nil
}
