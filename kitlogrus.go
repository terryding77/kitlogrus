package kitlogrus

import (
	"io"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/sirupsen/logrus"
	"github.com/terryding77/kitlogrus/formatter"
	"github.com/terryding77/kitlogrus/logger"
)

// Logger is the fundamental interface for all log operations. Log creates a
// log event from keyvals, a variadic sequence of alternating keys and values.
// Implementations must be safe for concurrent use by multiple goroutines. In
// particular, any implementation of Logger that appends to keyvals or
// modifies or retains any of its elements must make a copy first.
type Logger interface {
	Log(keyvals ...interface{}) error

	Print(args ...interface{}) // with level NO_LEVEL
	Debug(args ...interface{}) // with level DEBUG
	Info(args ...interface{})  // with level INFO
	Warn(args ...interface{})  // with level WARN
	Error(args ...interface{}) // with level ERROR

	Printf(format string, args ...interface{}) // with level NO_LEVEL
	Debugf(format string, args ...interface{}) // with level DEBUG
	Infof(format string, args ...interface{})  // with level INFO
	Warnf(format string, args ...interface{})  // with level WARN
	Errorf(format string, args ...interface{}) // with level ERROR

	KitLog() log.Logger     // Convert to github.com/go-kit/kit/log Logger struct
	Logrus() *logrus.Logger // Convert to github.com/go-kit/kit/log Logger struct
}

// NewJSONLogger return a new json logger
func NewJSONLogger(w io.Writer) Logger {
	return logger.New(&formatter.JSONFormatter{
		TimestampFormat:  time.RFC3339,
		DisableTimestamp: false,
		DataKey:          "data",
		PrettyPrint:      true,
	}, w, logrus.InfoLevel)
}

// NewLogfmtLogger return a new logfmt logger
func NewLogfmtLogger(w io.Writer) Logger {
	return logger.New(&formatter.TextFormatter{
		// Set to true to bypass checking for a TTY before outputting colors.
		ForceColors: false,
		// Force disabling colors.
		DisableColors:    false,
		TimestampFormat:  time.RFC3339,
		DisableTimestamp: false,
		// Disables the truncation of the level text to 4 characters.
		DisableLevelTruncation: false,
	}, w, logrus.InfoLevel)
}
