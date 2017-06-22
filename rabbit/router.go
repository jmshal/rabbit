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

	method := strings.ToUpper(r.Method)
	if method == "" {
		method = "GET"
	}

	for index, route := range a.config.Routes {
		for _, entrypoint := range route.Entrypoints {
			if len(entrypoint.Methods) > 0 {
				allow := false
				for _, verb := range entrypoint.Methods {
					if strings.ToUpper(verb) == method {
						allow = true
						break
					}
				}
				if !allow {
					continue
				}
			}
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
