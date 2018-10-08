package kitlogrus_test

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/terryding77/kitlogrus"
)

func logContainsTest(logFunc func(...interface{}), buf *bytes.Buffer, want string, args ...interface{}) {
	logFunc(args...)
	So(buf.String(), ShouldContainSubstring, want)
	buf.Reset()
}

func TestLog(t *testing.T) {
	buf := &bytes.Buffer{}
	Convey("testing", t, func() {
		Convey("json log", func() {
			log := kitlogrus.NewJSONLogger(buf)
			So(log, ShouldNotBeNil)

			logContainsTest(log.Debug, buf, "\"level\": \"debug\"", "hello", "world")
			logContainsTest(log.Info, buf, "\"level\": \"info\"", "hello", "world")
			logContainsTest(log.Warn, buf, "\"level\": \"warn\"", "hello", "world")
			logContainsTest(log.Error, buf, "\"level\": \"error\"", "hello", "world")

			logContainsTest(log.Debug, buf, "\"hello\": \"world\"", "hello", "world")
		})
		Convey("logfmt log", func() {
			log := kitlogrus.NewLogfmtLogger(buf)
			So(log, ShouldNotBeNil)

			logContainsTest(log.Debug, buf, "level=debug", "hello", "world")
			logContainsTest(log.Info, buf, "level=info", "hello", "world")
			logContainsTest(log.Warn, buf, "level=warn", "hello", "world")
			logContainsTest(log.Error, buf, "level=error", "hello", "world")

			logContainsTest(log.Debug, buf, "hello=world", "hello", "world")
		})
	})
}
