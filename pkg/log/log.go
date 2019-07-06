package log

import golog "log"

// Logger is for logging application messages
type Logger interface {
	Fatal(v ...interface{})
	Fatalln(v ...interface{})
	Fatalf(format string, v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type appLog struct {
}

//NewLogger provides the application logger
func NewLogger() Logger {
	return &appLog{}
}

func (appLog) Fatal(v ...interface{}) {
	golog.Fatal(v...)
}

func (appLog) Fatalln(v ...interface{}) {
	golog.Fatalln(v...)
}

func (appLog) Fatalf(format string, v ...interface{}) {
	golog.Fatalf(format, v...)
}

func (appLog) Printf(format string, v ...interface{}) {
	golog.Printf(format, v...)
}

func (appLog) Println(v ...interface{}) {
	golog.Println(v...)
}
