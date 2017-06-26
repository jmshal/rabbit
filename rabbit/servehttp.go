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
		var status int
		switch err {
		case ErrorNoMatchingRoute, ErrorWebsocketsNotAllowed:
			status = http.StatusBadRequest
		case ErrorNoRoutes:
			status = http.StatusInternalServerError
		case ErrorMethodNotAllowed:
			status = http.StatusMethodNotAllowed
		default:
			defer bugsnag.Notify(err, r)
			status = http.StatusInternalServerError
		}

		a.Log("%v <- %v %v (%v)", info.ID, r.Method, info.URL, err)
		http.Error(w, http.StatusText(status), status)
		return
	}
	info.Match = match

	if info.NeedsSecureRedirect() {
		a.SecureRedirect(w, r)
		return
	}

	a.ModifyRequest(r)
	info.ProxyURL = getRequestURL(r)

	var next func(http.ResponseWriter, *http.Request)
	if info.Websocket {
		next = a.wsProxy.ServeHTTP
	} else {
		next = a.httpProxy.ServeHTTP
	}

	a.Log("%v >> %v %v", info.ID, r.Method, info.ProxyURL)
	defer a.Log("%v <- %v %v", info.ID, r.Method, info.URL)

	next(w, r)
}
