package rabbit

import (
	"net/http"
	"net/url"

	"github.com/jmshal/wsutil"
)

// Builds a url.URL from a http.Request. Intentionally omits Userinfo.
func getRequestURL(r *http.Request) *url.URL {
	var s string
	isWebsocket := wsutil.IsWebSocketRequest(r)
	var isSecure bool
	if r.URL.Scheme != "" {
		isSecure = r.URL.Scheme == "https" || r.URL.Scheme == "wss"
	} else {
		isSecure = r.TLS != nil
	}
	if isSecure {
		if isWebsocket {
			s += "wss://"
		} else {
			s += "https://"
		}
	} else {
		if isWebsocket {
			s += "ws://"
		} else {
			s += "http://"
		}
	}
	if r.URL.Host != "" {
		s += r.URL.Host
	} else {
		s += r.Host
	}
	s += r.URL.Path
	if r.URL.RawQuery != "" {
		s += "?" + r.URL.RawQuery
	}
	url, _ := url.Parse(s)
	return url
}
