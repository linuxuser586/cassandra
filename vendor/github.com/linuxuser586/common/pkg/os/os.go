package os

import "os"

// Getenv retrieves the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present.
// To distinguish between an empty value and an unset value, use LookupEnv.
var Getenv = os.Getenv

// Hostname returns the host name reported by the kernel.
var Hostname = os.Hostname

// Exit causes the current program to exit with the given status code.
// Conventionally, code zero indicates success, non-zero an error.
// The program terminates immediately; deferred functions are not run.
var Exit = os.Exit

// Link creates newname as a hard link to the oldname file.
// If there is an error, it will be of type *LinkError.
var Link = os.Link
