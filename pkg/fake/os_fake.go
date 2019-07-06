package fake

import "github.com/linuxuser586/cassandra/pkg/virtual"

// OS is the fake os package.
type OS struct {
	Code int
}

// Exec is the fake exec package
type Exec struct {
	f string
}

// NewOS creates a new os package backed by fakes.
func NewOS() *OS {
	return &OS{Code: -1}
}

// NewExec creates a new exec package backed by fakes.
func NewExec(file string) virtual.Exec {
	return &Exec{f: file}
}

func (o *OS) Exit(code int) {
	o.Code = code
	panic(o)
}

func (e Exec) Command(name string, arg ...string) virtual.Cmd {
	return NewExecCmd(e.f)
}

// FailExec is used for testing exec package with expected failures.
type FailExec struct {
}

// NewFailExec creates an Exec packed use for for testing failures.
func NewFailExec() virtual.Exec {
	return &FailExec{}
}

func (e FailExec) Command(name string, arg ...string) virtual.Cmd {
	return NewFailExecCmd()
}
