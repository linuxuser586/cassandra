package virtual

import "os"

// OS is the virtual interface for os package.
type OS interface {
	Exit(code int)
}

// Exec is the virtual interface for exec package
type Exec interface {
	Command(name string, arg ...string) Cmd
}

type osPackage struct {
}

type execPackage struct {
}

// NewOS provides the system backing os package functions
func NewOS() OS {
	return &osPackage{}
}

// NewExec provides the system backing exec package functions
func NewExec() Exec {
	return &execPackage{}
}

// Exit is backed by os.Exit
func (p osPackage) Exit(code int) {
	os.Exit(code)
}

// Command is backed by exec.Command
func (p execPackage) Command(name string, arg ...string) Cmd {
	return NewExecCmd(name, arg...)
}
