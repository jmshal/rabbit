package rabbit

import (
	"net/http"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

func (a *rabbit) ServeHTTP(w http.ResponseWriter, _r *http.Request) {
	r, info := TagRequestInfo(_r)
	defer bugsnag.AutoNotify(r)

	a.Log("%v -> %v %v", info.ID, r.Method, info.URL)

	match, err := a.MatchRoute(r)
	if err != nil {
		defer bugsnag.Notify(err, r)

		var status int
		switch err {
		case ErrorNoMatchingRoute, ErrorWebsocketsNotAllowed:
			status = http.StatusBadRequest
		case ErrorNoRoutes:
			status = http.StatusInternalServerError
		case ErrorMethodNotAllowed:
			status = http.StatusMethodNotAllowed
		default:
			status = http.StatusInternalServerError
		}
		http.Error(w, http.StatusText(status), status)
		return
	}
	info.Match = match

	if match.Entrypoint.Secure && !info.Secure {
		a.SecureRedirect(w, r)
		return
	}

	a.ModifyRequest(r)

	var next func(http.ResponseWriter, *http.Request)
	if info.Websocket {
		next = a.wsProxy.ServeHTTP
	} else {
		next = a.httpProxy.ServeHTTP
	}

	fwdURL := getRequestURL(r)
	a.Log("%v <> %v %v", info.ID, r.Method, fwdURL)
	defer a.Log("%v <- %v %v", info.ID, r.Method, fwdURL)

	next(w, r)
}
