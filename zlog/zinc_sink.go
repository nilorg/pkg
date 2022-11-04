package zlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"go.uber.org/zap"
)

type ZincSink struct {
	username, password, host, index, format string
	httpClient                              *http.Client
}

// NewZincSink zinc://username:password@localhost:4080/index?format=json
func NewZincSink(uri *url.URL) (zs zap.Sink, err error) {
	var username, password string
	username = uri.User.Username()
	password, _ = uri.User.Password()
	zs = &ZincSink{
		username:   username,
		password:   password,
		host:       uri.Host,
		index:      uri.Path,
		format:     uri.Query().Get("format"),
		httpClient: new(http.Client),
	}
	return
}

func (zs *ZincSink) Close() (err error) {
	return
}

func (zs *ZincSink) Sync() (err error) {
	return
}

func (zs *ZincSink) Write(p []byte) (n int, err error) {
	zincURL := fmt.Sprintf("http://%s/%s", zs.host, path.Join("api", zs.index, "document"))
	var buf []byte
	if zs.format != "json" {
		if buf, err = json.Marshal(map[string]interface{}{
			"msg": strings.Fields(string(p)),
		}); err != nil {
			return
		}
	} else {
		buf = p
	}
	var req *http.Request
	if req, err = http.NewRequest("PUT", zincURL, bytes.NewBuffer(buf)); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(zs.username, zs.password)
	var resp *http.Response
	if resp, err = zs.httpClient.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()
	if code := resp.StatusCode >= 400; code {
		err = fmt.Errorf("%s", resp.Status)
		return
	}
	n = len(buf)
	return
}
