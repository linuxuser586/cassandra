package config

import (
	"fmt"
	"strings"
)

// RPCServerType are the options for the RPC Server.
type RPCServerType int

const (
	// Sync is one thread per thrift connection.
	Sync RPCServerType = iota
	// HSHA stands for "half synchronous, half asynchronous."
	HSHA
)

var rpcServerTypeID = map[RPCServerType]string{
	Sync: "sync",
	HSHA: "hsha",
}

var rpcServerTypeName = map[string]RPCServerType{
	"sync": Sync,
	"hsha": HSHA,
}

// MarshalYAML converts the enum to the string value for YAML.
func (r RPCServerType) MarshalYAML() (interface{}, error) {
	return rpcServerTypeID[r], nil
}

// UnmarshalYAML converts converts the YAML string to the enum integer.
func (r *RPCServerType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	if _, ok := rpcServerTypeName[s]; !ok {
		return fmt.Errorf("%v is not a valid RPC server type. Only sync, or hsha is valid", s)
	}
	*r = rpcServerTypeName[s]
	return nil
}
