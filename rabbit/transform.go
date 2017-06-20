package rabbit

import (
	"net/http"
	"path"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

const (
	requestIDHeader = "X-Request-Id"
)

var (
	removeHeaders = []string{
		"Server",
		"X-Powered-By",
		"X-AspNet-Version",
	}
)

func (a *Rabbit) TransformRequest(r *http.Request, entrypoint entrypoint, endpoint endpoint) error {
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
	r.URL.Path = path.Clean("/" + endpoint.Path + strippedPath)

	if requestID := r.Header.Get(requestIDHeader); requestID == "" {
		requestID = uuid.NewV4().String()
		r.Header.Set(requestIDHeader, requestID)
	}

	return nil
}

func (a *Rabbit) TransformResponse(r *http.Response) error {
	requestID := r.Request.Header.Get(requestIDHeader)
	for _, header := range removeHeaders {
		r.Header.Del(header)
	}
	r.Header.Set("Server", "rabbit")
	r.Header.Set(requestIDHeader, requestID)
	return nil
}
