package logger

type Logger interface {
	Panic(interface{})
	Fatal(interface{})
	Error(interface{})
	Warn(interface{})
	Info(interface{})
	Debug(interface{})
	Trace(interface{})
}
