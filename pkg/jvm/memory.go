package jvm

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const factor = 100
const jvmYG = "JVM_YG"
const jvmMem = "JVM_MEM"
const minMem = 1024
const maxMem = 8192

var cpus = runtime.NumCPU()
var mem uint64

func init() {
	c, err := ioutil.ReadFile("/sys/fs/cgroup/memory/memory.limit_in_bytes")
	if err != nil {
		log.Fatalf("could not get system memory: %v\n", err)
	}
	m, err := strconv.Atoi(strings.TrimSpace(string(c)))
	if err != nil {
		log.Fatalf("could not get system memory: %v\n", err)
	}
	// set value in MiB
	mem = uint64(m) / 1024.0 / 1024.0
}

// Xmn is the young generation memory.
func Xmn() string {
	x := "-Xmn"
	m := os.Getenv(jvmYG)
	if m != "" {
		return x + m
	}
	s := cpus * factor
	yg := s
	m = memory("")
	m = m[0 : len(m)-1]
	i, err := strconv.Atoi(m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not calculate young generation, using sensible: %v\n", err)
	} else {
		d := i / 4
		yg = d
		if d > s {
			yg = s
		}
	}
	o := fmt.Sprintf("%v%vm", x, yg)
	fmt.Printf("setting young generation: %v\n", o)
	return o
}

// Xms is the minimum heap.
func Xms() string {
	o := memory("-Xms")
	fmt.Printf("setting min memory: %v\n", o)
	return o
}

// Xmx is the maximum heap.
func Xmx() string {
	o := memory("-Xmx")
	fmt.Printf("setting max memory: %v\n", o)
	return o
}

func memory(x string) string {
	m := os.Getenv(jvmMem)
	if m != "" {
		return x + m
	}
	h := mem / 2
	q := mem / 4
	if h > minMem {
		h = minMem
	}
	if q > maxMem {
		q = maxMem
	}
	var meg uint64
	if h > q {
		meg = h
	} else {
		meg = q
	}
	return fmt.Sprintf("%v%vm", x, meg)
}
