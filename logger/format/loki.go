package format

import "github.com/sirupsen/logrus"

// LogrusLokiFormatter represents a loki format.
// It has logrus.Formatter which formats the entry and logrus.Fields which
// are added to the JSON message if not given in the entry data.
//
// Note: use the `DefaultFormatter` function to set a default loki formatter.
type LogrusLokiFormatter struct {
	logrus.Formatter
	logrus.Fields
}

var (
	lokiFields = logrus.Fields{"version": "1", "type": "log"}
)

// DefaultLogrusLokiFormatter returns a default loki formatter:
// A JSON format with "@version" set to "1" (unless set differently in `fields`,
// "type" to "log" (unless set differently in `fields`),
// "@timestamp" to the log time and "message" to the log message.
//
// Note: to set a different configuration use the `LogrusLokiFormatter` structure.
func DefaultLogrusLokiFormatter(fields logrus.Fields) logrus.Formatter {
	for k, v := range lokiFields {
		if _, ok := fields[k]; !ok {
			fields[k] = v
		}
	}

	return &LogrusLokiFormatter{
		Formatter: &logrus.TextFormatter{},
		Fields:    fields,
	}
}

// Format formats an entry to a loki format according to the given Formatter and Fields.
//
// Note: the given entry is copied and not changed during the formatting process.
func (f *LogrusLokiFormatter) Format(e *logrus.Entry) ([]byte, error) {
	ne := copyEntry(e, f.Fields)
	dataBytes, err := f.Formatter.Format(ne)
	releaseEntry(ne)
	return dataBytes, err
}
