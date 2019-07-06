package commit

import (
	"fmt"
	"strings"
)

// FailurePolicy is the numeric representation for the commit failure policy name.
type FailurePolicy int

const (
	// Die is to shut down the node and kill the JVM, so the node can be replaced.
	Die FailurePolicy = iota
	// Stop is to shut down the node, leaving the node effectively dead,
	// but can still be inspected via JMX.
	Stop
	// StopCommit is to shut down the commit log, letting writes collect but
	// continuing to service reads, as in pre-2.0.5 Cassandra
	StopCommit
	// Ignore fatal errors and let the batches fail
	Ignore
)

var failurePolicyTypeID = map[FailurePolicy]string{
	Die:        "die",
	Stop:       "stop",
	StopCommit: "stop_commit",
	Ignore:     "ignore",
}

var failurePolicyTypeName = map[string]FailurePolicy{
	"die":         Die,
	"stop":        Stop,
	"stop_commit": StopCommit,
	"ignore":      Ignore,
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
		return fmt.Errorf("%v is not a valid commit failure policy. Only die, stop, stop_commit, or ignore is valid", s)
	}
	*d = failurePolicyTypeName[s]
	return nil
}
