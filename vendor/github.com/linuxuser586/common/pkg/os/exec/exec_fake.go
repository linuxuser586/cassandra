package exec

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/linuxuser586/common/pkg/os/exec/mock"
)

var mockCmd *mock.MockCmd
var mockCmdCtx *mock.MockCmd

func fakeCommand(name string, arg ...string) Cmd {
	return mockCmd
}

func fakeCommandContext(ctx context.Context, name string, arg ...string) Cmd {
	return mockCmdCtx
}

// FakeCommand is backed by a mock and used for testing
func FakeCommand(t *testing.T) (*mock.MockCmd, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	execCommand = fakeCommand
	mockCmd = mock.NewMockCmd(ctrl)
	return mockCmd, ctrl
}

// FakeCommandContext is backed by a mock and used for testing
func FakeCommandContext(t *testing.T) (*mock.MockCmd, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	execCommandContext = fakeCommandContext
	mockCmdCtx = mock.NewMockCmd(ctrl)
	return mockCmdCtx, ctrl
}
