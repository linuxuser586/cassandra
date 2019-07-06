package v1

import (
	backup "github.com/linuxuser586/apis/grpc/cassandra/backup/v1"
	"github.com/spf13/afero"
)

func commitlog(in *backup.Downstream, stream backup.BackupService_StreamFromServer) error {
	files, err := afero.ReadDir(fs, commitPath)
	if err != nil {
		return err
	}
	bufSize := chunkSize(in.GetChunkSize())
	for _, fi := range files {
		if err := process(commitPath, fi, nil, bufSize, in, stream); err != nil {
			return err
		}
	}
	return nil
}

func restoreCommitlog(ds *backup.Downstream, m *backup.Upstream_Meta, c *backup.Upstream_Chunk) *backup.RestoreResponse {
	p := commitPath + m.GetFileName()
	return writeFile(p, ds, m, c)
}
