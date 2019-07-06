package net

import "net"

// InterfaceAddrs returns a list of the system's unicast interface
// addresses.
//
// The returned list does not identify the associated interface; use
// Interfaces and Interface.Addrs for more detail.
var InterfaceAddrs = net.InterfaceAddrs
