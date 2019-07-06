package commit

import (
	"fmt"
	"strings"
)

// LogSync is how the commit log should be synced
type LogSync int

const (
	// Periodic runs the sync in intervals
	Periodic LogSync = iota
	// Batch won’t ack writes until the commit log has been fsynced to disk.
	// It will wait commitlog_sync_batch_window_in_ms milliseconds between fsyncs.
	// This window should be kept short because the writer threads will be unable
	// to do extra work while waiting. (You may need to increase concurrent_writes
	// for the same reason.)
	Batch
)

var logSyncTypeID = map[LogSync]string{
	Periodic: "periodic",
	Batch:    "batch",
}

var logSyncTypeName = map[string]LogSync{
	"periodic": Periodic,
	"batch":    Batch,
}

// MarshalYAML converts the enum to the string value for YAML.
func (c LogSync) MarshalYAML() (interface{}, error) {
	return logSyncTypeID[c], nil
}

// UnmarshalYAML converts converts the YAML string to the enum integer.
func (c *LogSync) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	if _, ok := logSyncTypeName[s]; !ok {
		return fmt.Errorf("%v is not a valid commit log sync. Only periodic, or batch is valid", s)
	}
	*c = logSyncTypeName[s]
	return nil
}
