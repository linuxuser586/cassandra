package archive

import (
	"github.com/linuxuser586/cassandra/pkg/singleton/filesystem"
)

var fs = filesystem.Singleton()
