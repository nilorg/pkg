package logger

import (
	"testing"

	raven "github.com/getsentry/raven-go"
)

var (
	dsn = ""
)

func TestLog(t *testing.T) {
	tags := map[string]string{
		"server": "test api",
	}

	client, err := raven.NewClient(dsn, tags)
	if err != nil {
		t.Fatal(err)
	}

	err = SetHookSentry(client)
	if err != nil {
		t.Errorf("Logger RegisterSentry Failed: %v\n", err)
		return
	}
}
