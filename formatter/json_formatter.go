package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// Default key names for the default fields
const (
	FieldKeyMsg   = "msg"
	FieldKeyLevel = "level"
	FieldKeyTime  = "time"

	FieldValueLevel = "NONE"
)

var defaultTimestampFormat = time.RFC3339

// JSONFormatter formats logs into parsable json
type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// DataKey allows users to put all the log entry parameters into a nested dictionary at a given key.
	DataKey string

	// PrettyPrint will indent all json logs
	PrettyPrint bool
}

// Format renders a single log entry
func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+3)
	levelValue := FieldValueLevel
	for k, v := range entry.Data {
		if k == FieldKeyLevel {
			levelValue = v.(string)
		}
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	if f.DataKey != "" {
		newData := make(logrus.Fields, 4)
		newData[f.DataKey] = data
		data = newData
	}

	if f.TimestampFormat == "" {
		f.TimestampFormat = defaultTimestampFormat
	}

	if !f.DisableTimestamp {
		data[FieldKeyTime] = entry.Time.Format(f.TimestampFormat)
	}

	if entry.Message != "" {
		data[FieldKeyMsg] = entry.Message
	}

	if levelValue != FieldValueLevel {
		data[FieldKeyLevel] = levelValue
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	encoder := json.NewEncoder(b)
	if f.PrettyPrint {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}

	return b.Bytes(), nil
}
