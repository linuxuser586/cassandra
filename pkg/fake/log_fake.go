package fake

import (
	logger "github.com/linuxuser586/cassandra/pkg/log"
)

type fakeLog struct {
}

// NewLogger creates a new fake logger
func NewLogger() logger.Logger {
	return &fakeLog{}
}

func (fakeLog) Fatal(v ...interface{}) {
	panic(v[0])
}

func (fakeLog) Fatalln(v ...interface{}) {
	panic(v[0])
}

func (fakeLog) Fatalf(format string, v ...interface{}) {
	panic(v[0])
}

func (fakeLog) Printf(format string, v ...interface{}) {
	//TODO(Barry) implement
}

func (fakeLog) Println(v ...interface{}) {
	//TODO(Barry) implement
}
