package rabbit

import (
	"net/http"
)

type Rabbit struct {
	config *config
	mux    *http.ServeMux
}

func NewRabbit(config *config) *Rabbit {
	app := &Rabbit{
		config: config,
		mux:    http.NewServeMux(),
	}
	app.config.configureBugsnag()
	app.mux.HandleFunc("/", app.handleRequest)
	return app
}
