package exec

import (
	"context"
	"io"
	"os/exec"
	"syscall"
)

var execCommand = func(name string, arg ...string) Cmd {
	return &cmd{c: exec.Command(name, arg...)}
}

var execCommandContext = func(ctx context.Context, name string, arg ...string) Cmd {
	return &cmd{c: exec.CommandContext(ctx, name, arg...)}
}

// Cmd represents an external command being prepared or run.
//
// A Cmd cannot be reused after calling its Run, Output or CombinedOutput
// methods.
type Cmd interface {
	CombinedOutput() ([]byte, error)
	Credential(*syscall.Credential)
	Output() ([]byte, error)
	Run() error
	Start() error
	StderrPipe() (io.ReadCloser, error)
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	Wait() error
}

type cmd struct {
	c *exec.Cmd
}

// Command returns the Cmd struct to execute the named program with
// the given arguments.
//
// It sets only the Path and Args in the returned structure.
//
// If name contains no path separators, Command uses LookPath to
// resolve name to a complete path if possible. Otherwise it uses name
// directly as Path.
//
// The returned Cmd's Args field is constructed from the command name
// followed by the elements of arg, so arg should not include the
// command name itself. For example, Command("echo", "hello").
// Args[0] is always name, not the possibly resolved Path.
func Command(name string, arg ...string) Cmd {
	return execCommand(name, arg...)
}

// CommandContext is like Command but includes a context.
//
// The provided context is used to kill the process (by calling
// os.Process.Kill) if the context becomes done before the command
// completes on its own.
func CommandContext(ctx context.Context, name string, arg ...string) Cmd {
	return execCommandContext(ctx, name, arg...)
}

func (c *cmd) CombinedOutput() ([]byte, error) {
	return c.c.CombinedOutput()
}

func (c *cmd) Credential(cred *syscall.Credential) {
	c.c.SysProcAttr = &syscall.SysProcAttr{}
	c.c.SysProcAttr.Credential = cred
}

func (c *cmd) Output() ([]byte, error) {
	return c.c.Output()
}

func (c *cmd) Run() error {
	return c.c.Run()
}

func (c *cmd) Start() error {
	return c.c.Start()
}

func (c *cmd) StderrPipe() (io.ReadCloser, error) {
	return c.c.StderrPipe()
}

func (c *cmd) StdinPipe() (io.WriteCloser, error) {
	return c.c.StdinPipe()
}

func (c *cmd) StdoutPipe() (io.ReadCloser, error) {
	return c.c.StdoutPipe()
}

func (c *cmd) Wait() error {
	return c.c.Wait()
}
