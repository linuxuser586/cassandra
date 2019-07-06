package virtual

import (
	"os"
	"os/exec"
	"syscall"
)

// Cmd is the virtual structure for exec.Cmd
type Cmd interface {
	Output() ([]byte, error)
	Credential(*syscall.Credential)
	Stdout(*os.File)
	Stderr(*os.File)
	CombinedOutput() ([]byte, error)
	Run() error
	Process() *os.Process
}

// ExecCmd is for mapping to exec.Cmd.
type ExecCmd struct {
	cmd *exec.Cmd
}

// NewExecCmd provides the exec.Cmd backed by the system exec.Command
func NewExecCmd(name string, arg ...string) Cmd {
	c := exec.Command(name, arg...)
	return &ExecCmd{cmd: c}
}

// Output is backed by exec.Cmd.Output
func (c *ExecCmd) Output() ([]byte, error) {
	return c.cmd.Output()
}

// Credential is the system credential
func (c *ExecCmd) Credential(cred *syscall.Credential) {
	c.cmd.SysProcAttr = &syscall.SysProcAttr{}
	c.cmd.SysProcAttr.Credential = cred
}

// Stdout is standard out
func (c *ExecCmd) Stdout(f *os.File) {
	c.cmd.Stdout = f
}

// Stderr is standard error
func (c *ExecCmd) Stderr(f *os.File) {
	c.cmd.Stderr = f
}

func (c *ExecCmd) CombinedOutput() ([]byte, error) {
	return c.cmd.CombinedOutput()
}

// Run the command
func (c *ExecCmd) Run() error {
	return c.cmd.Run()
}

// Process is the OS process
func (c *ExecCmd) Process() *os.Process {
	return c.cmd.Process
}
