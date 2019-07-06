package config

import (
	"fmt"
	"strings"
)

// RowCacheClassName is the row cache implementation class name.
type RowCacheClassName int

const (
	// OHCProvider is the fully off-heap row cache implementation (default).
	OHCProvider RowCacheClassName = iota
	// SerializingCacheProvider is the row cache implementation availabile
	// in previous releases of Cassandra.
	SerializingCacheProvider
)

var rowCacheClassNameTypeID = map[RowCacheClassName]string{
	OHCProvider:              "org.apache.cassandra.cache.OHCProvider",
	SerializingCacheProvider: "org.apache.cassandra.cache.SerializingCacheProvider",
}

var rowCacheClassNameTypeName = map[string]RowCacheClassName{
	"org.apache.cassandra.cache.OHCProvider":              OHCProvider,
	"org.apache.cassandra.cache.SerializingCacheProvider": SerializingCacheProvider,
}

// MarshalYAML converts the enum to the string value for YAML.
func (r RowCacheClassName) MarshalYAML() (interface{}, error) {
	return rowCacheClassNameTypeID[r], nil
}

// UnmarshalYAML converts converts the YAML string to the enum integer.
func (r *RowCacheClassName) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	if _, ok := rowCacheClassNameTypeName[s]; !ok {
		return fmt.Errorf("%v is not a valid row cache class name. Only org.apache.cassandra.cache.OHCProvider, or org.apache.cassandra.cache.SerializingCacheProvider is valid", s)
	}
	*r = rowCacheClassNameTypeName[s]
	return nil
}
