package rabbit

import (
	"log"
	"net/http"
)

type Rabbit struct {
	config *config
	mux    *http.ServeMux
	logger *log.Logger
}

func NewRabbit(config *config) *Rabbit {
	app := &Rabbit{
		config: config,
		mux:    http.NewServeMux(),
	}
	app.setupLogger()
	app.Logger().Printf("configured %v TLS cert(s), and %v route(s)",
		len(app.config.Certs), len(app.config.Routes))
	app.configureBugsnag()
	app.mux.HandleFunc("/", app.handleRequest)
	return app
}
