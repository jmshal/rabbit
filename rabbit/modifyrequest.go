package rabbit

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	RequestID       = "X-Request-Id"
	XForwardedProto = "X-Forwarded-Proto"
	XForwardedPort  = "X-Forwarded-Port"
)

func (a *rabbit) ModifyRequest(r *http.Request) error {
	info := GetRequestInfo(r)
	endpoint := info.Match.Endpoint

	r.Header.Set(RequestID, info.ID)
	r.Header.Add(XForwardedProto, info.URL.Scheme)
	r.Header.Add(XForwardedPort, info.URL.Port())

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

	// TODO below, sort out retaining path without allowing path hijacking
	path := "/"
	path += strings.TrimLeft(endpoint.Path, "/")
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	path += strings.TrimLeft(strings.TrimPrefix(r.URL.Path, info.Match.Entrypoint.Path), "/")
	r.URL.Path = path

	return nil
}
