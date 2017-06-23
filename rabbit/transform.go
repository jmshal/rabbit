package rabbit

import (
	"net"
	"net/http"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

const (
	requestIDHeader = "X-Request-Id"
	xForwardedProto = "X-Forwarded-Proto"
	xForwardedPort  = "X-Forwarded-Port"
)

var (
	removeHeaders = []string{
		"Server",
		"Via",
		"X-Powered-By",
		"X-AspNet-Version",
	}
)

func (a *Rabbit) TransformRequest(r *http.Request, entrypoint entrypoint, endpoint endpoint) error {
	if r.TLS != nil {
		r.Header.Add(xForwardedProto, "https")
	} else {
		r.Header.Add(xForwardedProto, "http")
	}

	if _, port, err := net.SplitHostPort(r.Host); err == nil {
		r.Header.Add(xForwardedPort, port)
	}

	if endpoint.Origin != "" {
		r.Host = endpoint.Origin
	}

	if endpoint.Protocol == "https" {
		r.URL.Scheme = "https"
	} else {
		r.URL.Scheme = "http"
	}

	r.URL.Host = endpoint.Host
	if endpoint.Port != 0 {
		r.URL.Host += ":" + strconv.Itoa(int(endpoint.Port))
	}

	strippedPath := "/" + strings.TrimPrefix(r.URL.Path, entrypoint.Path)
	r.URL.Path = endpoint.Path + strippedPath

	if requestID := r.Header.Get(requestIDHeader); requestID == "" {
		requestID = uuid.NewV4().String()
		r.Header.Set(requestIDHeader, requestID)
	}

	a.Logger().Printf("-> %v %v %v", r.Method, r.Host, r.URL.Path)
	return nil
}

func (a *Rabbit) TransformResponse(r *http.Response) error {
	for _, header := range removeHeaders {
		r.Header.Del(header)
	}
	return nil
}
