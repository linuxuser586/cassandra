package v1

import (
	"encoding/base64"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"strings"

	"github.com/spf13/afero"

	backup "github.com/linuxuser586/apis/grpc/cassandra/backup/v1"
)

const (
	castagnoli = 0x82F63B78
	b1k        = 1024
	b2k        = 2048
	b4k        = 4096
	b8k        = 8192
	data       = "/data/"
)

var (
	dataLen = len(data)
)

func process(path string, fi os.FileInfo, err error, bufSize int, in *backup.Downstream, stream backup.BackupService_StreamFromServer) error {
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return nil
	}
	isDir, err := afero.IsDir(fs, path)
	if err != nil {
		return err
	}
	file := path
	if isDir {
		file = path + fi.Name()
	}
	hv, err := hashValue(in.GetHashType(), file)
	if err != nil {
		return err
	}
	relPath := ""
	if strings.Contains(path, data) {
		relPath = path[strings.Index(path, data)+dataLen : strings.LastIndex(path, "/")+1]
	}
	up := &backup.Upstream_Meta_{
		Meta: &backup.Upstream_Meta{
			FileName: fi.Name(),
			FileSize: fi.Size(),
			Hash:     hv,
			Path:     relPath,
		},
	}
	md := &backup.Upstream{Payload: up}
	if err := stream.Send(md); err != nil {
		return err
	}
	f, err := fs.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	pos := int64(0)
	for {
		buf := make([]byte, bufSize)
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n > 0 {
			c := &backup.Upstream_Chunk_{
				Chunk: &backup.Upstream_Chunk{
					Data:     buf[:n],
					Position: pos,
				},
			}
			chunk := &backup.Upstream{Payload: c}
			if err := stream.Send(chunk); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		pos += int64(n)
	}
	return nil
}

func chunkSize(cs backup.Downstream_ChunkSize) int {
	switch cs {
	case backup.Downstream_ONE_K:
		return b1k
	case backup.Downstream_TWO_K:
		return b2k
	case backup.Downstream_EIGHT_K:
		return b8k
	default:
		return b4k
	}
}

func hashValue(ht backup.Downstream_HashType, path string) (string, error) {
	if ht == backup.Downstream_CRC32C {
		f, err := fs.Open(path)
		if err != nil {
			return "", err
		}
		defer f.Close()
		t := crc32.MakeTable(castagnoli)
		h := crc32.New(t)
		if _, err := io.Copy(h, f); err != nil {
			return "", err
		}
		b := h.Sum(nil)[:]
		return base64.StdEncoding.EncodeToString(b), nil
	}
	return "", errors.New("unsupported hashtype")
}

func writeFile(p string, ds *backup.Downstream, m *backup.Upstream_Meta, c *backup.Upstream_Chunk) *backup.RestoreResponse {
	f, err := fs.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return &backup.RestoreResponse{
			Fail:    true,
			Message: err.Error(),
		}
	}
	_, err = f.WriteAt(c.GetData(), c.GetPosition())
	if err != nil {
		f.Close()
		return &backup.RestoreResponse{
			Fail:    true,
			Message: err.Error(),
		}
	}
	f.Close()
	if c.GetEof() {
		return checkFile(p, ds, m)
	}
	return &backup.RestoreResponse{Fail: false}
}

func checkFile(p string, ds *backup.Downstream, m *backup.Upstream_Meta) *backup.RestoreResponse {
	hv, err := hashValue(ds.GetHashType(), p)
	if err != nil {
		return &backup.RestoreResponse{
			Fail:    true,
			Message: err.Error(),
		}
	}
	if hv != m.GetHash() {
		return &backup.RestoreResponse{
			Fail:    true,
			Message: fmt.Sprintf("hash fail: got %s, expected %s", hv, m.GetHash()),
		}
	}
	return &backup.RestoreResponse{Fail: false}
}
