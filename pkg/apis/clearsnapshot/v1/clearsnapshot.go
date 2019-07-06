package v1

import (
	"context"
	"strings"
	"syscall"

	cs "github.com/linuxuser586/apis/grpc/cassandra/clearsnapshot/v1"
	nodetool "github.com/linuxuser586/cassandra/pkg/apis/nodetool/v1"
	"github.com/linuxuser586/cassandra/pkg/user"
	"github.com/linuxuser586/common/pkg/os/exec"
	"google.golang.org/grpc"
)

type service struct{}

// Register the GRPC server
func Register(s *grpc.Server) {
	cs.RegisterClearSnapshotServiceServer(s, &service{})
}

func (s *service) Run(ctx context.Context, in *cs.Request) (*cs.Response, error) {
	cmd := exec.Command("nodetool", args(in)...)
	cmd.Credential(&syscall.Credential{Uid: uint32(user.UID), Gid: uint32(user.GID)})
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &cs.Response{Message: string(out)}, nil
}

func args(in *cs.Request) []string {
	var s []string
	s = append(s, "clearsnapshot")
	if in.Args != nil {
		nodetool.ParseArgs(s, in.Args)
	}
	if in.SnapshotName != "" {
		s = append(s, "-t")
		s = append(s, in.SnapshotName)
	}
	if in.Keyspaces != "" {
		s = append(s, "--")
		for _, ks := range strings.Fields(in.Keyspaces) {
			s = append(s, ks)
		}
	}
	return s
}
