package disk

import (
	"fmt"
	"strings"
)

// OptimizationStrategy is the disk optimization to use.
type OptimizationStrategy int

const (
	// Spinning disk
	Spinning OptimizationStrategy = iota
	// SSD disk
	SSD
)

var optimizationStrategyTypeID = map[OptimizationStrategy]string{
	Spinning: "spinning",
	SSD:      "ssd",
}

var optimizationStrategyTypeName = map[string]OptimizationStrategy{
	"spinning": Spinning,
	"ssd":      SSD,
}

// MarshalYAML converts the enum to the string value for YAML.
func (o OptimizationStrategy) MarshalYAML() (interface{}, error) {
	return optimizationStrategyTypeID[o], nil
}

// UnmarshalYAML converts converts the YAML string to the enum integer.
func (o *OptimizationStrategy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	if _, ok := optimizationStrategyTypeName[s]; !ok {
		return fmt.Errorf("%v is not a valid disk optimization strategy. Only spinning, or ssd is valid", s)
	}
	*o = optimizationStrategyTypeName[s]
	return nil
}
