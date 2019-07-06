package logger

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"go.uber.org/zap"
)

var once sync.Once
var zapLogger *zap.Logger

// Zap is the single instance for the zap logger
func Zap() *zap.Logger {
	once.Do(func() {
		l, err := zap.NewProduction()
		if err != nil {
			fmt.Printf("failed to create logger: %v\n", err)
			os.Exit(1)
		}
		flush(l)
		zapLogger = l

	})
	return zapLogger
}

func flush(z *zap.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for {
			<-c
			z.Sync()
		}
	}()
}
