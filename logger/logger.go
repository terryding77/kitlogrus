package logger

import (
	"fmt"
	"io"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/sirupsen/logrus"
)

type logger struct {
	*logrus.Logger
}

var (
	printMessageKey = "msg"
	levelKey        = level.Key()
	levelError      = level.ErrorValue().String()
	levelWarn       = level.WarnValue().String()
	levelInfo       = level.InfoValue().String()
	levelDebug      = level.DebugValue().String()
)

func (l *logger) Log(keyvals ...interface{}) error {
	fields := logrus.Fields{}
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			fields[fmt.Sprint(keyvals[i])] = keyvals[i+1]
		} else {
			fields[fmt.Sprint(keyvals[i])] = log.ErrMissingValue
		}
	}
	l.WithFields(fields).Info()
	return nil
}

func (l *logger) Print(args ...interface{}) {
	l.Log(args...)
}

func (l *logger) printWithLevel(levelValue interface{}, args ...interface{}) {
	n := 2 + len(args)
	kvs := make([]interface{}, 0, n)
	kvs = append(kvs, levelKey)
	kvs = append(kvs, levelValue)
	kvs = append(kvs, args...)
	l.Log(kvs...)
}

func (l *logger) Debug(args ...interface{}) {
	l.printWithLevel(levelDebug, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.printWithLevel(levelInfo, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.printWithLevel(levelWarn, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.printWithLevel(levelError, args...)
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.Log(printMessageKey, fmt.Sprintf(format, args...))
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.Debug(printMessageKey, fmt.Sprintf(format, args...))
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.Info(printMessageKey, fmt.Sprintf(format, args...))
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.Warn(printMessageKey, fmt.Sprintf(format, args...))
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.Error(printMessageKey, fmt.Sprintf(format, args...))
}

func (l *logger) KitLog() log.Logger {
	return l
}

func (l *logger) Logrus() *logrus.Logger {
	return l.Logger
}

// New function return a logger struct for kitlogrus
func New(fmt logrus.Formatter, output io.Writer, logLevel logrus.Level) *logger {
	l := logrus.New()
	l.SetFormatter(fmt)
	l.SetOutput(output)
	l.SetLevel(logLevel)
	return &logger{Logger: l}
}
