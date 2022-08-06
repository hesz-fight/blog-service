package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// 日志公共字段
type Fields map[string]interface{}

// 日志等级
type Level int8

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

func (l *Level) String() string {
	switch *l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warm"
	case LevelError:
		return "error"
	case LevelFatal:
		return "Fatal"
	case LevelPanic:
		return "panic"
	}

	return ""
}

type Logger struct {
	logLogger *log.Logger
	ctx       context.Context
	fields    Fields
	callers   []string
}

func NewLogger(w io.Writer, prefix string, flag int) *Logger {
	l := log.New(w, prefix, flag) // 返回 log.Logger 实例
	return &Logger{logLogger: l}
}

// 返回和原本的 Logger 实例内容相同的新的实例
func (l *Logger) clone() *Logger {
	newLoger := *l

	return &newLoger
}

// 设置日志的公共字段并返回新的日志实例
func (l *Logger) WithFields(f Fields) *Logger {
	nl := l.clone()
	if nl.fields == nil {
		nl.fields = make(Fields)
	}

	for k, v := range f {
		nl.fields[k] = v
	}

	return nl
}

// 设置日志的上下文属性并返回新的日志实例
func (l *Logger) WithContext(ctx context.Context) *Logger {
	nl := l.clone()
	nl.ctx = ctx

	return nl
}

// 设置日志的调用栈信息并返回新的日志实例
func (l *Logger) WithCaller(skip int) *Logger {
	nl := l.clone()
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		f := runtime.FuncForPC(pc)
		nl.callers = []string{fmt.Sprintf("%s: %d %s", file, line, f.Name())}
	}

	return nl
}

// 设置当前整个调用栈的信息
func (l *Logger) WithCallersFrames() *Logger {

	minCallerDepth := 1
	maxCallerDepth := 25
	callers := []string{}
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	frame, more := frames.Next()
	for more {
		s := fmt.Sprintf("%s: %d %s", frame.File, frame.Line, frame.Function)
		callers = append(callers, s)
		frame, more = frames.Next()
	}

	nl := l.clone()
	nl.callers = callers

	return nl
}

// 日志格式化输出
// 将Logger实例的数据封装成键值对的形式返回
// 返回值用于输出显示
func (l *Logger) JSONFormat(level Level, message string) map[string]interface{} {

	data := make(Fields, len(l.fields)+4)

	data["level"] = level.String()
	data["time"] = time.Now().Local().UnixNano()
	data["message"] = message
	data["callers"] = l.callers

	if len(l.fields) > 0 { // 判空
		for k, v := range l.fields {
			if _, ok := data[k]; !ok { // 不存在
				data[k] = v
			}
		}
	}

	return data
}

// 追踪
func (l *Logger) WithTrace() *Logger {
	ginCtx, ok := l.ctx.(*gin.Context)
	if ok {
		return l.WithFields(Fields{
			"trace_id": ginCtx.MustGet("X-Trace-ID"),
			"span_id":  ginCtx.MustGet("X-Span-ID"),
		})
	}
	return l
}

// 日志输出方法
func (l *Logger) Output(level Level, message string) {

	body, _ := json.Marshal(l.JSONFormat(level, message))
	context := string(body)

	switch level {
	case LevelDebug, LevelInfo, LevelWarn, LevelError:
		l.logLogger.Print(context)
	case LevelFatal:
		l.logLogger.Fatal(context)
	case LevelPanic:
		l.logLogger.Panic(context)
	}
}

func (l *Logger) Debug(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelDebug, fmt.Sprint(v...))
}

func (l *Logger) Debugf(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelDebug, fmt.Sprintf(format, v...))
}

func (l *Logger) Info(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelInfo, fmt.Sprint(v...))
}

func (l *Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) InfofT(format string, v ...interface{}) {
	l.Output(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelWarn, fmt.Sprint(v...))
}

func (l *Logger) Warnf(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelWarn, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelError, fmt.Sprint(v...))
}

func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelError, fmt.Sprintf(format, v...))
}

func (l *Logger) ErrorfT(format string, v ...interface{}) {
	l.Output(LevelError, fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelFatal, fmt.Sprint(v...))
}

func (l *Logger) Fatalf(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelFatal, fmt.Sprintf(format, v...))
}

func (l *Logger) Panic(ctx context.Context, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelPanic, fmt.Sprint(v...))
}

func (l *Logger) Panicf(ctx context.Context, format string, v ...interface{}) {
	l.WithContext(ctx).WithTrace().Output(LevelPanic, fmt.Sprintf(format, v...))
}
