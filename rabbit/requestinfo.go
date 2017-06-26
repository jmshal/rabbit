package rabbit

import (
	"context"
	"net/http"
	"net/url"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	RequestInfoKey = "rabbit.RequestInfo"
)

type RequestInfo struct {
	ID             string
	URL            *url.URL
	ProxyURL       *url.URL
	Websocket      bool
	Secure         bool
	Match          *RouteMatch
	CachePolicy    *CachePolicy
	StartTime      time.Time
	Duration       time.Duration
	ResponseStatus int
}

func (i *RequestInfo) NeedsSecureRedirect() bool {
	return i.Match.Entrypoint.Secure && !i.Secure
}

func TagRequestInfo(r *http.Request) (*http.Request, *RequestInfo) {
	startTime := time.Now()
	var id string
	// if id = r.Header.Get(RequestID); id == "" {
	id = uuid.NewV4().String()
	// }
	u := getRequestURL(r)
	m := &RequestInfo{
		ID:        id,
		URL:       u,
		Websocket: u.Scheme == "wss" || u.Scheme == "ws",
		Secure:    r.TLS != nil,
		StartTime: startTime,
	}
	c := context.WithValue(r.Context(), RequestInfoKey, m)
	return r.WithContext(c), m
}

func GetRequestInfo(r *http.Request) *RequestInfo {
	return r.Context().Value(RequestInfoKey).(*RequestInfo)
}
