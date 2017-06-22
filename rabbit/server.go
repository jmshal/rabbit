package rabbit

import (
	"net/http"
	"net/http/httputil"

	"github.com/jmshal/wsutil"
)

func (a *Rabbit) handleRequest(w http.ResponseWriter, r *http.Request) {
	match, err := a.FindRoute(r)
	if err != nil {
		a.Logger().Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	endpoint, err := a.NextEndpoint(match.Route)
	if err != nil {
		a.Logger().Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var proxy http.Handler

	if wsutil.IsWebSocketRequest(r) {
		proxy = &wsutil.ReverseProxy{
			Director: func(r *http.Request) {
				a.TransformRequest(r, match.Entrypoint, *endpoint)
			},
			// TODO somehow add support for ModifyResponse
		}
	} else {
		proxy = &httputil.ReverseProxy{
			Director: func(r *http.Request) {
				a.TransformRequest(r, match.Entrypoint, *endpoint)
			},
			ModifyResponse: func(res *http.Response) error {
				return a.TransformResponse(res)
			},
		}
	}

	a.Logger().Printf("<- %v %v %v", r.Method, r.Host, r.URL.Path)
	proxy.ServeHTTP(w, r)
}
