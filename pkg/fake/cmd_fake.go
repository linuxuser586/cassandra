package fake

import (
	"io/ioutil"
	"os"
	"syscall"

	"github.com/linuxuser586/cassandra/pkg/virtual"
)

// Msg is the fake message
const Msg = "fake error"

// ErrFakeFatal is the fake fatal error.
type ErrFakeFatal struct {
}

func (f ErrFakeFatal) Error() string {
	return Msg
}

// ExecCmd is for mapping to exec.Cmd.
type execCmd struct {
	f string
}

// NewExecCmd creates a new fake exec cmd
func NewExecCmd(file string) virtual.Cmd {
	return &execCmd{f: file}
}

func (c *execCmd) Output() ([]byte, error) {
	b, err := ioutil.ReadFile(c.f)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *execCmd) Credential(cred *syscall.Credential) {

}

func (c *execCmd) Stdout(f *os.File) {

}

func (c *execCmd) Stderr(f *os.File) {

}

func (c *execCmd) CombinedOutput() ([]byte, error) {
	return []byte(""), nil
}

func (c *execCmd) Process() *os.Process {
	return &os.Process{}
}

func (c *execCmd) Run() error {
	return nil
}

type failExecCmd struct {
}

// NewFailExecCmd is used for testing a cmd that should fail.
func NewFailExecCmd() virtual.Cmd {
	return &failExecCmd{}
}

func (c failExecCmd) Output() ([]byte, error) {
	return nil, &ErrFakeFatal{}
}

func (c *failExecCmd) Credential(cred *syscall.Credential) {

}

func (c *failExecCmd) Stdout(f *os.File) {

}

func (c *failExecCmd) Stderr(f *os.File) {

}

func (c *failExecCmd) CombinedOutput() ([]byte, error) {
	return []byte(""), nil
}

func (c *failExecCmd) Process() *os.Process {
	return &os.Process{}
}

func (c *failExecCmd) Run() error {
	return nil
}
