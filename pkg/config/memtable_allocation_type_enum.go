package config

import (
	"fmt"
	"strings"
)

// MemtableAllocationType is for memtable memory allocation.
type MemtableAllocationType int

const (
	// HeapBuffers are on heap nio buffers.
	HeapBuffers MemtableAllocationType = iota
	// OffHeapBuffers are off heap (direct) nio buffers
	OffHeapBuffers
	// OffHeapObjects are off heap objects
	OffHeapObjects
)

var memtableAllocationTypeID = map[MemtableAllocationType]string{
	HeapBuffers:    "heap_buffers",
	OffHeapBuffers: "offheap_buffers",
	OffHeapObjects: "offheap_objects",
}

var memtableAllocationTypeName = map[string]MemtableAllocationType{
	"heap_buffers":    HeapBuffers,
	"offheap_buffers": OffHeapBuffers,
	"offheap_objects": OffHeapObjects,
}

// MarshalYAML converts the enum to the string value for YAML.
func (m MemtableAllocationType) MarshalYAML() (interface{}, error) {
	return memtableAllocationTypeID[m], nil
}

// UnmarshalYAML converts converts the YAML string to the enum integer.
func (m *MemtableAllocationType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	if _, ok := memtableAllocationTypeName[s]; !ok {
		return fmt.Errorf("%v is not a valid memtable allocation type. Only heap_buffers, offheap_buffers, or offheap_objects is valid", s)
	}
	*m = memtableAllocationTypeName[s]
	return nil
}
