package readyprobe

import (
	"testing"

	"github.com/linuxuser586/cassandra/pkg/fake"
	"github.com/linuxuser586/common/pkg/os"
)

type testData struct {
	file string
	ip   string
	code int
}

func TestSingleNodeReady(t *testing.T) {
	d := &testData{"testdata/single-node-ready.txt", "172.17.0.2", 0}
	runCurrent(t, d)
}

func TestThreeNodesCurrentNodeNotReady(t *testing.T) {
	d := &testData{"testdata/three-nodes-current-node-not-ready.txt", "172.18.0.2", 1}
	runCurrent(t, d)
}

func TestThreeNodesCurrentNodeReady(t *testing.T) {
	d := &testData{"testdata/three-nodes-current-node-ready.txt", "172.19.0.2", 0}
	runCurrent(t, d)
}

func runCurrent(t *testing.T, tbl *testData) {
	os.Exit = func(code int) {

	}
	// fake exit panics and is captured here
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case *fake.OS:
				if v.Code != tbl.code {
					t.Fatalf("got: %v, want: exit status %v for test %v", v.Code, tbl.code, tbl.file)
				}
			default:
				t.Fatalf("%#v\n", v)
			}
		}
	}()
	log = fake.NewLogger()
	// Run()
}

func TestOutputError(t *testing.T) {
	// fake fail exec panics and is captured here
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case *fake.ErrFakeFatal:
				if v.Error() != fake.Msg {
					t.Fatalf("got: %v, want: %v ", v.Error(), fake.Msg)
				}
			default:
				t.Fatalf("%#v\n", v)
			}
		}
	}()
	log = fake.NewLogger()
	//Run()
}
