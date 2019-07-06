package v1

import nodetool "github.com/linuxuser586/apis/grpc/cassandra/nodetool/v1"

// ParseArgs for Nodetool
func ParseArgs(s []string, o *nodetool.Args) {
	if o.Host != "" {
		s = append(s, "--host")
		s = append(s, o.Host)
	}
	if o.Port != "" {
		s = append(s, "--port")
		s = append(s, o.Port)
	}
	if o.Username != "" {
		s = append(s, "--username")
		s = append(s, o.Host)
	}
	if o.Password != "" {
		s = append(s, "--password")
		s = append(s, o.Password)
	}
	if o.PasswordFile != "" {
		s = append(s, "--password-file")
		s = append(s, o.PasswordFile)
	}
}
