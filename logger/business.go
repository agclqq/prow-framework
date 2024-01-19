package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/SkyAPM/go2sky"

	"github.com/agclqq/prow-framework/skywalking"
	"github.com/agclqq/prow-framework/times"

	"github.com/sirupsen/logrus"
)

// BusinessLog 以下内容实现了接口
type BusinessLog struct {
	*logrus.Logger
	ctx       context.Context
	withLine  bool
	withTrace bool
	fields    []map[string]interface{}
	entry     *logrus.Entry
}

type LogLevel string

const (
	LogLevelPanic   LogLevel = "panic"
	LogLevelFatal            = "fatal"
	LogLevelError            = "error"
	LogLevelWarning          = "warning"
	LogLevelInfo             = "info"
	LogLevelDebug            = "debug"
	LogLevelTrace            = "trace"
)
const DefaultRetain = 7 //日志保留天数
const maximumCallerDepth = 26
const knownLogrusFrames int = 4

var minimumCallerDepth = 1
var currentPackage = "github.com/agclqq/prow-framework/logger"
var businessOnce sync.Once

type Option func(l *BusinessLog)

func New(opts ...Option) *BusinessLog {
	ln := logrus.New()
	ln.SetReportCaller(false) //启用后，请求定位不准
	ln.Formatter = &logrus.JSONFormatter{
		TimestampFormat:   times.FormatDatetimeMicro,
		DisableTimestamp:  false,
		DisableHTMLEscape: true,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier:  nil,
		PrettyPrint:       false,
	}
	ln.SetOutput(os.Stdout)
	lg := &BusinessLog{Logger: ln}
	for _, opt := range opts {
		opt(lg)
	}
	return lg
}
func WithContext(ctx context.Context) Option {
	return func(l *BusinessLog) {
		l.ctx = ctx
	}
}

func WithFile(file string, retain uint) Option {
	return func(l *BusinessLog) {
		if retain <= 0 {
			retain = DefaultRetain
		}
		writer, err := RotateDailyLog(file, retain)
		if err != nil {
			_ = fmt.Errorf("%s", err)
		}
		l.Logger.SetOutput(writer)
	}
}
func WithLevel(level LogLevel) Option {
	return func(l *BusinessLog) {
		if level, err := logrus.ParseLevel(string(level)); err != nil {
			l.Logger.SetLevel(level)
		}
	}
}

func WithLine(b bool) Option {
	return func(l *BusinessLog) {
		l.withLine = b
	}
}

func WithField(key string, value interface{}) Option {
	return func(l *BusinessLog) {
		l.fields = append(l.fields, map[string]interface{}{key: value})
	}
}

func WithTrace(b bool) Option {
	return func(l *BusinessLog) {
		l.withTrace = b
	}
}

//func (bs *BusinessLog) WithContext(ctx context.Context) *BusinessLog {
//	bs.ctx = ctx
//	return bs
//}
//
//func (bs *BusinessLog) WithFile(file string, retain uint) *BusinessLog {
//	if retain <= 0 {
//		retain = DefaultRetain
//	}
//	writer, err := RotateDailyLog(file, retain)
//	if err != nil {
//		_ = fmt.Errorf("%s", err)
//	}
//	bs.Logger.SetOutput(writer)
//	return bs
//}
//func (bs *BusinessLog) WithLevel(level LogLevel) {
//	if level, err := logrus.ParseLevel(string(level)); err != nil {
//		bs.Logger.SetLevel(level)
//	}
//}
//
//func (bs *BusinessLog) WithLine(b bool) *BusinessLog {
//	bs.withLine = b
//	return bs
//}
//func (bs *BusinessLog) WithField(key string, value interface{}) *BusinessLog {
//	bs.fields = append(bs.fields, map[string]interface{}{key: value})
//	return bs
//}
//func (bs *BusinessLog) WithTrace(b bool) *BusinessLog {
//	bs.withTrace = b
//	return bs
//}

//func NewStdoutLog() *BusinessLog {
//	ln := logrus.New()
//	ln.SetReportCaller(false) //启用后，请求定位不准
//	ln.Formatter = &logrus.JSONFormatter{
//		TimestampFormat:   times.FormatDatetimeMicro,
//		DisableTimestamp:  false,
//		DisableHTMLEscape: true,
//		DataKey:           "",
//		FieldMap:          nil,
//		CallerPrettyfier:  nil,
//		PrettyPrint:       false,
//	}
//	ln.SetOutput(os.Stdout)
//	BusinessLogger = &BusinessLog{Logger: ln}
//	return BusinessLogger
//}
//
//func NewBusinessLog() *BusinessLog {
//	ln := logrus.New()
//	ln.SetReportCaller(false) //启用后，请求定位不准
//	ln.Formatter = &logrus.JSONFormatter{
//		TimestampFormat:   times.FormatDatetimeMicro,
//		DisableTimestamp:  false,
//		DisableHTMLEscape: true,
//		DataKey:           "",
//		FieldMap:          nil,
//		CallerPrettyfier:  nil,
//		PrettyPrint:       false,
//	}
//	businessLog := config.GetLog("business")
//
//	//f, err := file.OpenOrCreate(businessLog["file"])
//	//if err != nil {
//	//	panic(err)
//	//}
//	retain, err := strconv.Atoi(businessLog["retain"])
//	if err != nil {
//		retain = DefaultRetain
//	}
//	writer, err := RotateDailyLog(businessLog["file"], uint(retain))
//	if err != nil {
//		_ = fmt.Errorf("%s", err)
//	}
//	ln.SetOutput(writer)
//	if level, err := logrus.ParseLevel(businessLog["level"]); err != nil {
//		ln.SetLevel(level)
//	}
//	BusinessLogger = &BusinessLog{Logger: ln}
//	return BusinessLogger
//}

func businessLogCaller() *runtime.Frame {
	// cache this package's fully-qualified name
	businessOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "businessLogCaller") {
				currentPackage = getPackageName(funcName)
				break
			}
		}

		minimumCallerDepth = knownLogrusFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != currentPackage {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}

func (bs *BusinessLog) Panic(v interface{}) {
	withFields(bs).Panic(v)
}
func (bs *BusinessLog) Panicf(format string, v ...interface{}) {
	withFields(bs).Panicf(format, v...)
}
func (bs *BusinessLog) Fatal(v interface{}) {
	withFields(bs).Fatalln(v)
}
func (bs *BusinessLog) Fatalf(format string, v ...interface{}) {
	withFields(bs).Fatalf(format, v...)
}
func (bs *BusinessLog) Error(v interface{}) {
	withFields(bs).Errorln(v)
}
func (bs *BusinessLog) Errorf(format string, v ...interface{}) {
	withFields(bs).Errorf(format, v...)
}
func (bs *BusinessLog) Warn(v interface{}) {
	withFields(bs).Warnln(v)
}
func (bs *BusinessLog) Warnf(format string, v ...interface{}) {
	withFields(bs).Warnf(format, v...)
}
func (bs *BusinessLog) Info(v interface{}) {
	withFields(bs).Infoln(v)
}
func (bs *BusinessLog) Infof(format string, v ...interface{}) {
	withFields(bs).Infof(format, v...)
}
func (bs *BusinessLog) Debug(v interface{}) {
	withFields(bs).Debugln(v)
}
func (bs *BusinessLog) Debugf(format string, v ...interface{}) {
	withFields(bs).Debugf(format, v...)
}
func (bs *BusinessLog) Trace(v interface{}) {
	withFields(bs).Traceln(v)
}
func (bs *BusinessLog) Tracef(format string, v ...interface{}) {
	withFields(bs).Tracef(format, v...)
}

func withFileLine(bs *BusinessLog) {
	if bs.withLine {
		frame := businessLogCaller()
		bs.entry.WithField("file", frame.File+":"+strconv.Itoa(frame.Line))
	}
}
func withTrace(bs *BusinessLog) {
	traceId := skywalking.GetTraceId(bs.ctx)
	if traceId != go2sky.EmptyTraceID {
		bs.entry.WithField("traceId", traceId)
	}

}
func withFields(bs *BusinessLog) *logrus.Entry {
	if bs.entry == nil {
		bs.entry = logrus.NewEntry(bs.Logger)
	}
	withFileLine(bs)
	withTrace(bs)
	for _, field := range bs.fields {
		for k, v := range field {
			bs.entry.WithField(k, v)
		}
	}
	return bs.entry
}
