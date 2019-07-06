package v1

import (
	"context"
	"io/ioutil"
	"testing"

	backup "github.com/linuxuser586/apis/grpc/cassandra/backup/v1"
	"github.com/spf13/afero"
	"google.golang.org/grpc/metadata"
)

const (
	logDir = "testdata/commitlog/"
	tmpDir = "/tmp/out/"
)

func TestCommitLog(t *testing.T) {
	fs = afero.NewMemMapFs()
	loadCommitLogs(t)
	ds := &backup.Downstream{
		Type: backup.Downstream_COMMITLOG,
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
		file string
		hash string
	}{
		{"CommitLog-3-129398221.log", "9dpuXg=="},
		{"CommitLog-3-129398222.log", "I+xJTQ=="},
		{"CommitLog-3-129398223.log", "0YfKTg=="},
		{"test.bin", "5oB10A=="},
	}
	for _, tt := range td {
		meta := &backup.Upstream_Meta{
			FileName: tt.file,
			Hash:     tt.hash,
		}
		verifyFile(t, dir, meta)
	}
}

func TestDeleteCommitLog(t *testing.T) {
	fs = afero.NewMemMapFs()
	loadCommitLogs(t)
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
		service.recv <- &backup.DeleteRequest{File: "CommitLog-3-129398221.log"}
		service.recv <- &backup.DeleteRequest{File: "data-does-not-exist.log"}
		service.recv <- &backup.DeleteRequest{File: "../CommitLog-3-129398222.log"}
		close(service.recv)
	}()
	s := backupService{}
	if err := s.DeleteCommitLog(service); err != nil {
		t.Error(err)
	}
	close(service.in)
	<-waitc
}

func loadCommitLogs(t *testing.T) {
	files, err := ioutil.ReadDir(logDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, fi := range files {
		rpath := logDir + fi.Name()
		wpath := commitPath + fi.Name()
		b, err := ioutil.ReadFile(rpath)
		if err != nil {
			t.Fatal(err)
		}
		if err := afero.WriteFile(fs, wpath, b, 644); err != nil {
			t.Fatal(err)
		}
		t.Logf("Added: %s", wpath)
	}
}

func recv(t *testing.T, dir afero.Fs, service *fakeService) error {
	var file afero.File
	for in := range service.in {
		if in.GetMeta() != nil {
			m := in.GetMeta()
			t.Logf("Meta: %#v", m)
			var err error
			file, err = dir.Create(m.GetPath() + m.GetFileName())
			if err != nil {
				return err
			}
		} else {
			c := in.GetChunk()
			t.Logf("Position: %v, Data Len: %v", c.GetPosition(), len(c.GetData()))
			t.Logf("File: %v", file.Name())
			if _, err := file.WriteAt(c.GetData(), c.GetPosition()); err != nil {
				return err
			}
		}
	}
	return nil
}

func verifyFile(t *testing.T, dir afero.Fs, meta *backup.Upstream_Meta) {
	path := tmpDir + meta.GetPath() + meta.GetFileName()
	hv, err := hashValue(backup.Downstream_CRC32C, path)
	if err != nil {
		t.Error(err)
		return
	}
	if hv != meta.GetHash() {
		t.Errorf("Hash: got: %v, want: %v", hv, meta.GetHash())
	}
}

func newFakeService() *fakeService {
	s := &fakeService{}
	s.in = make(chan *backup.Upstream, 20)
	return s
}

type fakeService struct {
	in chan *backup.Upstream
	fakeServerStream
}

func (f *fakeService) Send(in *backup.Upstream) error {
	f.in <- in
	return nil
}

type fakeServerStream struct{}

func (f *fakeServerStream) SetHeader(md metadata.MD) error {
	return nil
}

func (f *fakeServerStream) SendHeader(md metadata.MD) error {
	return nil
}

func (f *fakeServerStream) SetTrailer(md metadata.MD) {

}

func (f *fakeServerStream) Context() context.Context {
	return nil
}

func (f *fakeServerStream) SendMsg(m interface{}) error {
	return nil
}

func (f *fakeServerStream) RecvMsg(m interface{}) error {
	return nil
}
