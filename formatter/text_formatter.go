package formatter

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/sirupsen/logrus"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 36
	gray    = 37
)

var (
	baseTimestamp time.Time
	// emptyFieldMap FieldMap
)

func init() {
	baseTimestamp = time.Now()
}

// TextFormatter formats logs into text
type TextFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// Disables the truncation of the level text to 4 characters.
	DisableLevelTruncation bool

	// Whether the logger's out is to a terminal
	isTerminal bool

	terminalInitOnce sync.Once
}

func (f *TextFormatter) init(entry *logrus.Entry) {
	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)

		if f.isTerminal {
			initTerminal(entry.Logger.Out)
		}
	}
}

func (f *TextFormatter) isColored() bool {
	isColored := f.ForceColors || f.isTerminal

	return isColored && !f.DisableColors
}

// Format renders a single log entry
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	keys := make([]string, 0, len(entry.Data))
	levelValue := FieldValueLevel
	for k := range entry.Data {
		if k == FieldKeyLevel {
			levelValue = entry.Data[k].(string)
		}
		keys = append(keys, k)
	}

	fixedKeys := make([]string, 0, 3+len(entry.Data))
	if !f.DisableTimestamp {
		fixedKeys = append(fixedKeys, FieldKeyTime)
	}

	if levelValue != FieldValueLevel {
		fixedKeys = append(fixedKeys, FieldKeyLevel)
	}

	if _, ok := entry.Data[FieldKeyMsg]; entry.Message != "" && !ok {
		fixedKeys = append(fixedKeys, FieldKeyMsg)
	}

	fixedKeys = append(fixedKeys, keys...)

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.terminalInitOnce.Do(func() { f.init(entry) })

	if f.TimestampFormat == "" {
		f.TimestampFormat = defaultTimestampFormat
	}
	if f.isColored() {
		f.printColored(b, entry, keys, f.TimestampFormat, levelValue)
	} else {
		for _, key := range fixedKeys {
			var value interface{}
			switch key {
			case FieldKeyTime:
				value = entry.Time.Format(f.TimestampFormat)
			case FieldKeyLevel:
				value = levelValue
			case FieldKeyMsg:
				if entry.Message != "" {
					value = entry.Message
				} else {
					value = entry.Data[key]
				}
			default:
				value = entry.Data[key]
			}
			f.appendKeyValue(b, key, value)
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string, timestampFormat string, levelValue string) {
	levelText := strings.ToUpper(levelValue)
	if !f.DisableLevelTruncation {
		levelText = levelText[0:4]
	}

	var levelColor int
	switch levelValue {
	case level.ErrorValue().String():
		levelColor = red
	case level.WarnValue().String():
		levelColor = yellow
	case level.InfoValue().String():
		levelColor = green
	case level.DebugValue().String():
		levelColor = gray
	case FieldValueLevel:
		levelText = "NONE"
		levelColor = blue
	default:
		levelColor = blue
	}

	// Remove a single newline if it already exists in the message to keep
	// the behavior of logrus text_formatter the same as the stdlib log package
	entry.Message = strings.TrimSuffix(entry.Message, "\n")

	if f.DisableTimestamp {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m %-s ", levelColor, levelText, entry.Message)
	} else {
		fmt.Fprintf(b, "\x1b[%dm%s[%s]\x1b[0m %-s ", levelColor, levelText, entry.Time.Format(timestampFormat), entry.Message)
	}
	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, k)
		f.appendValue(b, v)
	}
}

func (f *TextFormatter) needsQuoting(text string) bool {
	if len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}
