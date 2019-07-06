package v1

import (
	"github.com/linuxuser586/cassandra/pkg/singleton/filesystem"
	"github.com/linuxuser586/common/pkg/logger"
)

var zaplog = logger.Zap()
var log = zaplog.Sugar()
var fs = filesystem.Singleton()
