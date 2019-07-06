package nodetool

import (
	"fmt"
	goos "os"
	"os/signal"
	"syscall"

	"github.com/linuxuser586/cassandra/pkg/jvm"
	"github.com/linuxuser586/cassandra/pkg/user"
	"github.com/linuxuser586/cassandra/pkg/virtual"
)

var stdout = goos.Stdout
var stderr = goos.Stderr
var exec = virtual.NewExec()

// Run nodetool
func Run() {
	cmd := exec.Command("java", nodetoolArgs()...)
	cmd.Credential(&syscall.Credential{Uid: uint32(user.UID), Gid: uint32(user.GID)})
	cmd.Stdout(stdout)
	cmd.Stderr(stderr)
	c := make(chan goos.Signal, 1)
	signal.Notify(c, goos.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println("stopping cassandra nodetool")
			cmd.Process().Signal(sig)
		}
	}()
	cmd.Run()
}

func nodetoolArgs() []string {
	s := jvm.Args()
	s = append(s, mem())
	s = append(s, "org.apache.cassandra.tools.NodeTool")
	s = append(s, "-p")
	s = append(s, "7199")
	a := goos.Args[1:]
	s = append(s, a...)
	return s
}

func mem() string {
	x := "-Xmx"
	m := goos.Getenv("NODETOOL_MEM")
	if m != "" {
		return x + m
	}
	return x + "128m"
}
