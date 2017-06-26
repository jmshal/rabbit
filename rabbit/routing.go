package rabbit

import (
	"errors"
	"net/http"
	"strings"

	glob "github.com/ryanuber/go-glob"
)

var (
	ErrorNoRoutes             = errors.New("no registered routes, config error")
	ErrorNoMatchingRoute      = errors.New("no matching routes for request")
	ErrorMethodNotAllowed     = errors.New("method not allowed")
	ErrorWebsocketsNotAllowed = errors.New("websocket request not allowed")
)

type RouteMatch struct {
	Route        *Route
	Entrypoint   *Entrypoint
	Endpoint     *Endpoint
	TrailingPath string
}

func (a *rabbit) MatchRoute(r *http.Request) (*RouteMatch, error) {
	if len(a.config.Routes) == 0 {
		return nil, ErrorNoRoutes
	}

	info := GetRequestInfo(r)

	for _, route := range a.config.Routes {
		for _, entrypoint := range route.Entrypoints {
			// check host (glob)
			if entrypoint.Host != "" &&
				!glob.Glob(entrypoint.Host, info.URL.Hostname()) {
				continue
			}

			// check path contains prefix
			if entrypoint.Path != "" &&
				!strings.HasPrefix(info.URL.Path, entrypoint.Path) {
				continue
			}

			// check supported method
			if len(entrypoint.Methods) > 0 {
				allow := false
				for _, verb := range entrypoint.Methods {
					if strings.ToUpper(verb) == r.Method {
						allow = true
						break
					}
				}
				if !allow {
					return nil, ErrorMethodNotAllowed
				}
			}

			// check if websockets are supported
			if info.Websocket && !entrypoint.Websockets {
				return nil, ErrorWebsocketsNotAllowed
			}

			trailingPath := strings.TrimPrefix(r.URL.Path, entrypoint.Path)
			trailingPath = "/" + strings.TrimLeft(trailingPath, "/")

			return &RouteMatch{
				Route:        route,
				Entrypoint:   entrypoint,
				Endpoint:     route.Endpoints[0], // TODO lb
				TrailingPath: trailingPath,
			}, nil
		}
	}

	return nil, ErrorNoMatchingRoute
}
