package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

var levelNames = map[LogLevel]string{
	LevelDebug:    "DEBUG",
	LevelInfo:     "INFO",
	LevelWarning:  "WARNING",
	LevelError:    "ERROR",
	LevelCritical: "CRITICAL",
}

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Critical(v ...interface{})
	Criticalf(format string, v ...interface{})
	SetLevel(level LogLevel)
	SetOutput(w io.Writer)
	WithFields(fields map[string]interface{}) Logger
}

type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	Fields    map[string]interface{}
	Caller    string
}

type BotLogger struct {
	level      LogLevel
	output     io.Writer
	mu         sync.Mutex
	timeFormat string
	caller     bool
	fields     map[string]interface{}
}

func NewLogger() *BotLogger {
	return &BotLogger{
		level:      LevelInfo,
		output:     os.Stdout,
		timeFormat: "2006-01-02 15:04:05",
		caller:     true,
		fields:     make(map[string]interface{}),
	}
}

func (l *BotLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *BotLogger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
}

func (l *BotLogger) WithFields(fields map[string]interface{}) Logger {
	newLogger := NewLogger()
	newLogger.level = l.level
	newLogger.output = l.output
	newLogger.timeFormat = l.timeFormat
	newLogger.caller = l.caller

	newLogger.fields = make(map[string]interface{})
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

func (l *BotLogger) log(level LogLevel, message string, v ...interface{}) {
	if level < l.level {
		return
	}

	var caller string
	if l.caller {
		_, file, line, ok := runtime.Caller(3)
		if ok {
			caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
		}
	}

	msg := fmt.Sprintf(message, v...)
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Caller:    caller,
		Fields:    l.fields,
	}

	output := fmt.Sprintf("[%s] %s %s",
		entry.Timestamp.Format(l.timeFormat),
		levelNames[level],
		entry.Message)

	if entry.Caller != "" {
		output += fmt.Sprintf(" (%s)", entry.Caller)
	}

	if len(entry.Fields) > 0 {
		var fields []string
		for k, v := range entry.Fields {
			fields = append(fields, fmt.Sprintf("%s=%v", k, v))
		}
		output += " " + strings.Join(fields, " ")
	}

	output += "\n"

	l.mu.Lock()
	defer l.mu.Unlock()
	_, err := fmt.Fprint(l.output, output)
	if err != nil {
		log.Printf("failed to write log: %v", err)
	}

	if level == LevelCritical {
		os.Exit(1)
	}
}

func (l *BotLogger) Debug(v ...interface{}) {
	l.log(LevelDebug, fmt.Sprint(v...))
}

func (l *BotLogger) Debugf(format string, v ...interface{}) {
	l.log(LevelDebug, format, v...)
}

func (l *BotLogger) Info(v ...interface{}) {
	l.log(LevelInfo, fmt.Sprint(v...))
}

func (l *BotLogger) Infof(format string, v ...interface{}) {
	l.log(LevelInfo, format, v...)
}

func (l *BotLogger) Warn(v ...interface{}) {
	l.log(LevelWarning, fmt.Sprint(v...))
}

func (l *BotLogger) Warnf(format string, v ...interface{}) {
	l.log(LevelWarning, format, v...)
}

func (l *BotLogger) Error(v ...interface{}) {
	l.log(LevelError, fmt.Sprint(v...))
}

func (l *BotLogger) Errorf(format string, v ...interface{}) {
	l.log(LevelError, format, v...)
}

func (l *BotLogger) Critical(v ...interface{}) {
	l.log(LevelCritical, fmt.Sprint(v...))
}

func (l *BotLogger) Criticalf(format string, v ...interface{}) {
	l.log(LevelCritical, format, v...)
}

var DefaultLogger Logger = NewLogger()

func Debug(v ...interface{})                    { DefaultLogger.Debug(v...) }
func Debugf(format string, v ...interface{})    { DefaultLogger.Debugf(format, v...) }
func Info(v ...interface{})                     { DefaultLogger.Info(v...) }
func Infof(format string, v ...interface{})     { DefaultLogger.Infof(format, v...) }
func Warn(v ...interface{})                     { DefaultLogger.Warn(v...) }
func Warnf(format string, v ...interface{})     { DefaultLogger.Warnf(format, v...) }
func Error(v ...interface{})                    { DefaultLogger.Error(v...) }
func Errorf(format string, v ...interface{})    { DefaultLogger.Errorf(format, v...) }
func Critical(v ...interface{})                 { DefaultLogger.Critical(v...) }
func Criticalf(format string, v ...interface{}) { DefaultLogger.Criticalf(format, v...) }
