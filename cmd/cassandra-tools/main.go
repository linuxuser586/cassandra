package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/linuxuser586/cassandra/cmd/archive"
	"github.com/linuxuser586/cassandra/cmd/cassandra"
	"github.com/linuxuser586/cassandra/cmd/nodetool"
	"github.com/linuxuser586/cassandra/cmd/readyprobe"
	"github.com/linuxuser586/cassandra/pkg/cert"
	"github.com/linuxuser586/cassandra/pkg/server/rpc"
)

const (
	bin = "/usr/local/bin/"
	src = bin + "cassandra"
	nt  = "nodetool"
	rp  = "ready-probe"
	cp  = "copy-commitlog"
	ln  = "link-commitlog"
)

func main() {
	app := os.Args[0]
	if strings.HasSuffix(app, "cassandra") {
		if err := link(); err != nil {
			log.Fatal(err)
		}
		if err := cert.Setup(); err != nil {
			log.Fatal(err)
		}
		go rpc.Start()
		cassandra.Start()
	} else if strings.HasSuffix(app, nt) {
		nodetool.Run()
	} else if strings.HasSuffix(app, rp) {
		readyprobe.Run()
	} else if strings.HasSuffix(app, cp) {
		s, d, err := archiveArgs()
		if err != nil {
			return
		}
		if err := archive.Copy(s, d); err != nil {
			log.Errorf("commitlog copy error: %v\n", err)
		}
	} else if strings.HasSuffix(app, ln) {
		s, d, err := archiveArgs()
		if err != nil {
			return
		}
		if err := archive.Link(s, d); err != nil {
			log.Errorf("commitlog link error: %v\n", err)
		}
	} else {
		fmt.Printf("%v command not found\n", app)
	}
}

func archiveArgs() (string, string, error) {
	err := false
	if len(os.Args) < 2 {
		log.Error("missing arg 1: source file")
		err = true
	}
	if len(os.Args) < 3 {
		log.Error("missing arg 2: destination file")
		err = true
	}
	if err {
		return "", "", errors.New("missing arg")
	}
	s := os.Args[1]
	d := os.Args[2]
	return s, d, nil
}

func link() error {
	dsts := []string{nt, rp, cp, ln}
	for _, dst := range dsts {
		if err := os.Link(src, bin+dst); err != nil {
			return err
		}
	}
	return nil
}
