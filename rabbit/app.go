package rabbit

import (
	"net/http"
	"sync"
)

type Rabbit struct {
	config *config
	mux    sync.Mutex
	server *http.ServeMux
}

func NewRabbit(config *config) *Rabbit {
	app := &Rabbit{
		config: config,
		server: http.NewServeMux(),
	}
	app.config.configureBugsnag()
	app.server.HandleFunc("/", app.handleRequest)
	return app
}
