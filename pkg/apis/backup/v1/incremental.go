package v1

import (
	"io"
	"os"
	"regexp"
	"strings"

	backup "github.com/linuxuser586/apis/grpc/cassandra/backup/v1"
	"github.com/spf13/afero"
)

const backupPathReFormat = "/backups/*(.)*$"

var backupPathRe = regexp.MustCompile(backupPathReFormat)

func incremental(in *backup.Downstream, stream backup.BackupService_StreamFromServer) error {
	bufSize := chunkSize(in.GetChunkSize())
	err := afero.Walk(fs, dataPath, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() || !strings.Contains(path, "/backups/") {
			return nil
		}
		return process(path, fi, err, bufSize, in, stream)
	})
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func restoreIncremental(ds *backup.Downstream, m *backup.Upstream_Meta, c *backup.Upstream_Chunk) *backup.RestoreResponse {
	p := dataPath + m.GetPath() + m.GetFileName()
	return writeFile(p, ds, m, c)
}
