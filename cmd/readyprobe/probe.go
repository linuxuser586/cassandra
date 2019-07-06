package readyprobe

import (
	"strings"

	logger "github.com/linuxuser586/cassandra/pkg/log"
	"github.com/linuxuser586/common/pkg/net"
	"github.com/linuxuser586/common/pkg/os"
	"github.com/linuxuser586/common/pkg/os/exec"
)

var log = logger.NewLogger()
var podIP = net.PodIP

// Run ready probe
func Run() {
	b, err := exec.Command("nodetool", "status").Output()
	if err != nil {
		log.Fatal(err)
	}
	ip, err := podIP()
	if err != nil {
		log.Fatal(err)
	}
	s := strings.Split(string(b), "\n")
	for _, l := range s {
		if strings.Contains(l, ip) && strings.HasPrefix(l, "UN") {
			// UN (Up Normal)
			// probe returns ready
			os.Exit(0)
		}
	}
	// probe returns not ready
	os.Exit(1)
}
