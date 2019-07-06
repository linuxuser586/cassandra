package bootstrap

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const mapCnt = "/proc/sys/vm/max_map_count"

func main() {
	setVM()
}

func setVM() {
	vm := []byte("1048575")
	err := ioutil.WriteFile(mapCnt, vm, 0644)
	if err != nil {
		fmt.Printf("could not adjust vm.max_map_count: %v\n", err)
		fmt.Println("for Docker, use --privileged with the run command ")
		log.Fatalln("for Kubernetes, set privileged true in the container securityContext capabilities section")
	}

	c, err := ioutil.ReadFile(mapCnt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read vm.max_map_count: %v\n", err)
	}
	i, err := strconv.Atoi(strings.TrimSpace((string(c))))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not convert vm.max_map_count to integer: %v\n", err)
	}
	if i < 1048575 {
		fmt.Printf("vm.max_map_count is %v and must be greater or equal to 1048575", i)
		fmt.Println("for Docker, use --privileged with the run command ")
		fmt.Println("or run sudo sysctl -w vm.max_map_count=1048575 on the host system")
		log.Fatalln("for Kubernetes, run the initContainer")
	} else {
		fmt.Printf("vm.max_map_count: %v", string(c))
	}
}
