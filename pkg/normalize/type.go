package normalize

import "strings"

// True determins is the string should be treated as true.
// It evaluates case insensitive true, yes, y, and 1 as true.
func True(s string) bool {
	s = strings.ToLower(s)
	return s == "true" || s == "yes" || s == "y" || s == "1"
}
