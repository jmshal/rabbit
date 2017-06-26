package rabbit

import (
	"fmt"
	"net/http"
	"os"
)

const (
	RequestID        = "X-Request-Id"
	XForwardedProto  = "X-Forwarded-Proto"
	XForwardedPort   = "X-Forwarded-Port"
	XForwardedServer = "X-Forwarded-Server"
)

var (
	hostname = ""
)

func (a *rabbit) ModifyRequest(r *http.Request) error {
	info := GetRequestInfo(r)
	endpoint := info.Match.Endpoint

	r.Header.Set(RequestID, info.ID)
	r.Header.Add(XForwardedProto, info.URL.Scheme)
	r.Header.Add(XForwardedPort, info.URL.Port())

	if hostname != "" {
		r.Header.Add(XForwardedServer, hostname)
	}

	if endpoint.Origin != "" {
		r.Host = endpoint.Origin
	}

	if endpoint.Scheme == "https" {
		if info.Websocket {
			r.URL.Scheme = "wss"
		} else {
			r.URL.Scheme = "https"
		}
	} else {
		if info.Websocket {
			r.URL.Scheme = "ws"
		} else {
			r.URL.Scheme = "http"
		}
	}

	r.URL.Host = endpoint.Host
	if endpoint.Port != 0 {
		r.URL.Host += fmt.Sprintf(":%v", endpoint.Port)
	}

	var path string
	if endpoint.Path != "" {
		path += endpoint.Path
	}
	if info.Match.TrailingPath != "/" {
		path += info.Match.TrailingPath
	}
	if path == "" {
		path = "/"
	}
	r.URL.Path = path

	return nil
}

func init() {
	hostname, _ = os.Hostname()
}
