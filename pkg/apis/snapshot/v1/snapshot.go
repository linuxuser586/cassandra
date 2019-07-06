package v1

import (
	"context"
	"log"
	"strings"
	"syscall"

	snapshot "github.com/linuxuser586/apis/grpc/cassandra/snapshot/v1"
	nodetool "github.com/linuxuser586/cassandra/pkg/apis/nodetool/v1"
	"github.com/linuxuser586/cassandra/pkg/user"
	"github.com/linuxuser586/cassandra/pkg/virtual"
	"google.golang.org/grpc"
)

var exec = virtual.NewExec()

type snapshotter struct{}

// Register the GRPC server
func Register(s *grpc.Server) {
	snapshot.RegisterSnapshotterServer(s, &snapshotter{})
}

func (s *snapshotter) Snapshot(ctx context.Context, in *snapshot.Request) (*snapshot.Response, error) {
	cmd := exec.Command("nodetool", args(in)...)
	cmd.Credential(&syscall.Credential{Uid: uint32(user.UID), Gid: uint32(user.GID)})
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &snapshot.Response{Code: 200, Message: string(out)}, nil
}

func args(in *snapshot.Request) []string {
	var s []string
	s = append(s, "snapshot")
	if in.Args != nil {
		nodetool.ParseArgs(s, in.Args)
	}
	if in.Table != "" {
		s = append(s, "--table")
		s = append(s, in.Table)
	}
	if in.KtList != "" {
		s = append(s, "--kt-list")
		s = append(s, in.KtList)
	}
	if in.SkipFlush {
		s = append(s, "--skip-flush")
	}
	if in.Tag != "" {
		s = append(s, "--tag")
		s = append(s, in.Tag)
	}
	if in.Keyspaces != "" {
		s = append(s, "--")
		for _, ks := range strings.Fields(in.Keyspaces) {
			s = append(s, ks)
		}
	}
	return s
}
