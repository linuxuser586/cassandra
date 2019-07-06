package config

import (
	"fmt"
	"strings"
)

// InternodeCompression controls whether traffic between nodes is compressed.
type InternodeCompression int

const (
	// All traffic is compressed.
	All InternodeCompression = iota
	// DC traffic between different datacenters is compressed.
	DC
	// None does not compress anything.
	None
)

var internodeCompressionTypeID = map[InternodeCompression]string{
	All:  "all",
	DC:   "dc",
	None: "none",
}

var internodeCompressionTypeName = map[string]InternodeCompression{
	"all":  All,
	"dc":   DC,
	"none": None,
}

// MarshalYAML converts the enum to the string value for YAML.
func (i InternodeCompression) MarshalYAML() (interface{}, error) {
	return internodeCompressionTypeID[i], nil
}

// UnmarshalYAML converts converts the YAML string to the enum integer.
func (i *InternodeCompression) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	if _, ok := internodeCompressionTypeName[s]; !ok {
		return fmt.Errorf("%v is not a valid internode compression. Only all, dc, or none is valid", s)
	}
	*i = internodeCompressionTypeName[s]
	return nil
}
