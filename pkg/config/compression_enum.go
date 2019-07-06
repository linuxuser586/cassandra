package config

import (
	"fmt"
	"strings"
)

// CompressionType is the numeric representation for compression name.
type CompressionType int

// CompressionParameters is a string map of strings
type CompressionParameters map[string]string

const (
	// DeflateCompressor compresses hints using the Deflate algorithm.
	DeflateCompressor CompressionType = iota
	// LZ4Compressor compresses using the LZ4 algorithm.
	LZ4Compressor
	// SnappyCompressor compresses hints using the Snappy algorithm.
	SnappyCompressor
)

var compressionTypeID = map[CompressionType]string{
	DeflateCompressor: "DeflateCompressor",
	LZ4Compressor:     "LZ4Compressor",
	SnappyCompressor:  "SnappyCompressor",
}

var compressionTypeName = map[string]CompressionType{
	"DeflateCompressor": DeflateCompressor,
	"LZ4Compressor":     LZ4Compressor,
	"SnappyCompressor":  SnappyCompressor,
}

// MarshalYAML converts the enum to the string value for YAML.
func (c CompressionType) MarshalYAML() (interface{}, error) {
	return compressionTypeID[c], nil
}

// UnmarshalYAML converts converts the YAML string to the enum integer.
func (c *CompressionType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	if _, ok := compressionTypeName[s]; !ok {
		return fmt.Errorf("%v is not a valid compressor. Only DeflateCompressor, LZ4Compressor, or SnappyCompressor is valid", s)
	}
	*c = compressionTypeName[s]
	return nil
}
