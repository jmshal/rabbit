package rabbit

import "errors"

var (
	ErrorNoEndpoint = errors.New("route does not have any endpoints")
)

func (a *Rabbit) NextEndpoint(route route) (*endpoint, error) {
	if len(route.Endpoints) == 0 {
		return nil, ErrorNoEndpoint
	}
	// TODO Allow for multiple endpoints (rr)
	return &route.Endpoints[0], nil
}
