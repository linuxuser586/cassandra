package disk

import (
	"fmt"
	"strings"
)

// FailurePolicy is the numeric representation for the disk failure policy name.
type FailurePolicy int

const (
	// Die is to shut down gossip and client transports and kill the JVM
	// for any fs errors or single-sstable errors, so the node can be replaced.
	Die FailurePolicy = iota
	// StopParanoid is to shut down gossip and client transports even for
	// single-sstable errors, kill the JVM for errors during startup.
	StopParanoid
	// Stop is to shut down gossip and client transports, leaving the node
	// effectively dead, but can still be inspected via JMX, kill the JVM
	// for errors during startup.
	Stop
	// BestEffort is to stop using the failed disk and respond to requests
	// based on remaining available sstables. This means you WILL see obsolete
	// data at CL.ONE!
	BestEffort
	// Ignore fatal errors and let requests fail, as in pre-1.2 Cassandra
	Ignore
)

var failurePolicyTypeID = map[FailurePolicy]string{
	Die:          "die",
	StopParanoid: "stop_paranoid",
	Stop:         "stop",
	BestEffort:   "best_effort",
	Ignore:       "ignore",
}

var failurePolicyTypeName = map[string]FailurePolicy{
	"die":           Die,
	"stop_paranoid": StopParanoid,
	"stop":          Stop,
	"best_effort":   BestEffort,
	"ignore":        Ignore,
}

// MarshalYAML converts the enum to the string value for YAML.
func (d FailurePolicy) MarshalYAML() (interface{}, error) {
	return failurePolicyTypeID[d], nil
}

// UnmarshalYAML converts converts the YAML string to the enum integer.
func (d *FailurePolicy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	if _, ok := failurePolicyTypeName[s]; !ok {
		return fmt.Errorf("%v is not a valid disk failure policy. Only die, stop_paranoid, stop, best_effort, or ignore is valid", s)
	}
	*d = failurePolicyTypeName[s]
	return nil
}
