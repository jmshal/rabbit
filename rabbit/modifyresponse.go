package rabbit

import (
	"net/http"
)

var (
	removeHeaders = []string{
		"Server",
		"Via",
		"X-Powered-By",
		"X-AspNet-Version",
	}
)

func (a *rabbit) modifyResponse(r *http.Response) error {
	info := GetRequestInfo(r.Request)
	a.Log("%v << %v %v", info.ID, r.Request.Method, info.ProxyURL)

	// extract any important details about the response
	info.ResponseStatus = r.StatusCode

	for _, header := range removeHeaders {
		r.Header.Del(header)
	}

	return nil
}
