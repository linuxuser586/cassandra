// Package v1 is the version 1 backup responsible for managing the raw files
package v1

import (
	"errors"
	"io"
	"strings"

	backup "github.com/linuxuser586/apis/grpc/cassandra/backup/v1"
	"google.golang.org/grpc"
)

const (
	basePath   = "/var/lib/cassandra/"
	commitPath = basePath + "backup/commitlog/"
	dataPath   = basePath + "data/"
)

type backupService struct{}

// Register the GRPC server
func Register(s *grpc.Server) {
	backup.RegisterBackupServiceServer(s, &backupService{})
}

func (b *backupService) StreamFrom(in *backup.Downstream, stream backup.BackupService_StreamFromServer) error {
	if in.GetType() == backup.Downstream_SNAPSHOT {
		return snapshot(in, stream)
	} else if in.GetType() == backup.Downstream_INCREMENTAL {
		return incremental(in, stream)
	}
	return commitlog(in, stream)
}

func (b *backupService) StreamTo(stream backup.BackupService_StreamToServer) error {
	var ds *backup.Downstream
	var m *backup.Upstream_Meta
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		n, ds := getDownstream(ds, in)
		if ds == nil {
			return errors.New("no downstream found")
		}
		if n {
			// next downstream found and now need to find the upstream
			continue
		}
		us := in.GetUpstream()
		if us == nil {
			return errors.New("no upstream found")
		}
		n, m := getMetadata(m, us)
		if m == nil {
			return errors.New("no metadata found")
		}
		if n {
			// next metadata found and now need to process the chunks
			continue
		}
		c := us.GetChunk()
		if c == nil {
			return errors.New("could not find chunk")
		}
		if err := stream.Send(restore(ds, m, c)); err != nil {
			return err
		}
	}
}

func getDownstream(cds *backup.Downstream, r *backup.Restore) (bool, *backup.Downstream) {
	ds := r.GetDownstream()
	if ds != nil {
		return true, ds
	}
	return false, cds
}

func getMetadata(cm *backup.Upstream_Meta, us *backup.Upstream) (bool, *backup.Upstream_Meta) {
	m := us.GetMeta()
	if m != nil {
		return true, m
	}
	return false, cm
}

func restore(ds *backup.Downstream, m *backup.Upstream_Meta, c *backup.Upstream_Chunk) *backup.RestoreResponse {
	if ds.GetType() == backup.Downstream_SNAPSHOT {
		return restoreSnapshot(ds, m, c)
	} else if ds.GetType() == backup.Downstream_INCREMENTAL {
		return restoreIncremental(ds, m, c)
	}
	return restoreCommitlog(ds, m, c)
}

func (b *backupService) DeleteCommitLog(stream backup.BackupService_DeleteCommitLogServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if strings.Contains(in.GetFile(), "/") {
			stream.Send(&backup.DeleteResponse{Fail: true, Message: "file must not contain any directories"})
			continue
		}
		f := commitPath + in.GetFile()
		if err := fs.Remove(f); err != nil {
			stream.Send(&backup.DeleteResponse{Fail: true, Message: err.Error()})
			continue
		}
		stream.Send(&backup.DeleteResponse{Fail: false})
	}
}

func (b *backupService) DeleteIncremental(stream backup.BackupService_DeleteIncrementalServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if !backupPathRe.MatchString(in.GetFile()) {
			stream.Send(&backup.DeleteResponse{Fail: true, Message: "file must match pattern " + backupPathReFormat})
			continue
		}
		f := dataPath + in.GetFile()
		if err := fs.Remove(f); err != nil {
			stream.Send(&backup.DeleteResponse{Fail: true, Message: err.Error()})
			continue
		}
		stream.Send(&backup.DeleteResponse{Fail: false})
	}
}
