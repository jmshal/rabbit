package rabbit

import (
	"errors"
	"net/http"
	"strings"

	glob "github.com/ryanuber/go-glob"
)

var (
	ErrorNoRoute = errors.New("no route found")
)

type RouteMatch struct {
	Route      route
	Entrypoint entrypoint
}

func (a *Rabbit) FindRoute(r *http.Request) (*RouteMatch, error) {
	host, port, err := parseHostPort(r.Host)

	if err != nil {
		return nil, err
	}

	a.mux.Lock()
	defer a.mux.Unlock()

	for _, route := range a.config.Routes {
		for _, entrypoint := range route.Entrypoints {
			if entrypoint.Host != "" &&
				!glob.Glob(entrypoint.Host, host) {
				continue
			}
			if entrypoint.Port != 0 &&
				entrypoint.Port != port {
				continue
			}
			if entrypoint.Path != "" &&
				!strings.HasPrefix(r.URL.Path, entrypoint.Path) {
				continue
			}
			return &RouteMatch{route, entrypoint}, nil
		}
	}

	return nil, ErrorNoRoute
}
