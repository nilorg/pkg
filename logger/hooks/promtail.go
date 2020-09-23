package hooks

import (
	"fmt"
	"time"

	"github.com/afiskon/promtail-client/promtail"
	"github.com/sirupsen/logrus"
)

// PromtailLogrusHook ...
type PromtailLogrusHook struct {
	Client    promtail.Client
	formatter logrus.Formatter
	LogLevels []logrus.Level
}

// NewPromtailLogrusHook ...
func NewPromtailLogrusHook(client promtail.Client, formatter logrus.Formatter, channel string) logrus.Hook {
	return &PromtailLogrusHook{
		Client:    client,
		formatter: formatter,
		LogLevels: logrus.AllLevels,
	}
}

// NewDefaultPromtailLogrusHook ...
func NewDefaultPromtailLogrusHook(hostURL, namespace string, formatter logrus.Formatter) logrus.Hook {
	labels := fmt.Sprintf(`{namespace="%s"}`, namespace)
	conf := promtail.ClientConfig{
		PushURL:            hostURL + "/api/prom/push",
		Labels:             labels,
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		SendLevel:          promtail.DEBUG,
	}
	client, _ := promtail.NewClientProto(conf)

	return &PromtailLogrusHook{
		Client:    client,
		formatter: formatter,
		LogLevels: logrus.AllLevels,
	}
}

// Fire ...
func (h *PromtailLogrusHook) Fire(e *logrus.Entry) error {
	dataBytes, err := h.formatter.Format(e)
	if err != nil {
		return err
	}

	switch e.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		h.Client.Debugf(string(dataBytes))
	case logrus.InfoLevel:
		h.Client.Infof(string(dataBytes))
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		h.Client.Errorf(string(dataBytes))
	default:
		h.Client.Warnf(string(dataBytes))
	}
	return nil
}

// Levels returns all logrus levels.
func (h *PromtailLogrusHook) Levels() []logrus.Level {
	return h.LogLevels
}
