package rpc

import (
	"sync"
	"testing"

	"go.uber.org/zap"
)

var once sync.Once

func setup(t *testing.T) {
	once.Do(func() {
		var err error
		zaplog, err = zap.NewDevelopment()
		if err != nil {
			t.Error(err)
		}
		log = zaplog.Sugar()
	})
}

func TestServer(t *testing.T) {
	setup(t)
	//Start()
}
