package rpc

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	backup "github.com/linuxuser586/cassandra/pkg/apis/backup/v1"
	cs "github.com/linuxuser586/cassandra/pkg/apis/clearsnapshot/v1"
	snapshot "github.com/linuxuser586/cassandra/pkg/apis/snapshot/v1"
	"github.com/linuxuser586/cassandra/release"
	"google.golang.org/grpc"
)

// Start the server
func Start() string {
	p := ":" + os.Getenv("CASSANDRA_GRPC_PORT")
	if p == ":" {
		p = ":10042"
	}
	log.Infof("Starting Cassandra %s GRPC %s-%s on %s", release.CassandraVersion, release.Version, release.Commit, p)
	lis, err := net.Listen("tcp", p)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	s := grpc.NewServer()
	backup.Register(s)
	cs.Register(s)
	snapshot.Register(s)
	go func() {
		for {
			<-c
			log.Info("Stopping Cassandra GRPC server")
			s.Stop()
		}
	}()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	return "stopped Cassandra GRPC server"
}
