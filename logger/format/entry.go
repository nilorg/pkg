package format

import (
	"sync"

	"github.com/nilorg/pkg/logger"
	"github.com/nilorg/sdk/log"
	"github.com/sirupsen/logrus"
)

// Using a pool to re-use of old entries when formatting Logstash messages.
// It is used in the Fire function.
var entryPool = sync.Pool{
	New: func() interface{} {
		return &logrus.Entry{}
	},
}

// copyEntry copies the entry `e` to a new entry and then adds all the fields in `fields` that are missing in the new entry data.
// It uses `entryPool` to re-use allocated entries.
func copyEntry(e *logrus.Entry, fields logrus.Fields) *logrus.Entry {
	ne := entryPool.Get().(*logrus.Entry)
	ne.Message = e.Message
	ne.Level = e.Level
	ne.Time = e.Time
	ne.Caller = e.Caller
	ne.Data = logrus.Fields{}
	if e.Context != nil {
		var (
			traceID string
			spanID  string
			userID  string
			ok      bool
		)
		if traceID, ok = log.FromTraceIDContext(e.Context); ok {
			ne.Data[logger.TraceIDKey] = traceID
		}
		if spanID, ok = log.FromSpanIDContext(e.Context); ok {
			ne.Data[logger.SpanIDKey] = spanID
		}
		if userID, ok = log.FromUserIDContext(e.Context); ok {
			ne.Data[logger.UserIDKey] = userID
		}
	}
	for k, v := range fields {
		ne.Data[k] = v
	}
	for k, v := range e.Data {
		ne.Data[k] = v
	}
	return ne
}

// releaseEntry puts the given entry back to `entryPool`. It must be called if copyEntry is called.
func releaseEntry(e *logrus.Entry) {
	entryPool.Put(e)
}
