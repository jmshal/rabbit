package rabbit

import (
	"crypto/tls"
	"net"
	"net/http"
	"strconv"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

func getPortHandler(a *Rabbit) (port string, handler http.Handler) {
	port = ":" + strconv.Itoa(int(a.config.Server.Port))
	handler = bugsnag.Handler(a.server)
	return
}

func (a *Rabbit) listenHTTP() error {
	port, handler := getPortHandler(a)
	return http.ListenAndServe(port, handler)
}

func (a *Rabbit) listenTLS() error {
	port, handler := getPortHandler(a)

	server := &http.Server{
		Addr:    port,
		Handler: handler,
	}

	config := &tls.Config{}
	if server.TLSConfig != nil {
		*config = *server.TLSConfig
	}

	if config.NextProtos == nil {
		config.NextProtos = []string{"h2", "h2-14", "http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, len(a.config.TLS))
	for index, cert := range a.config.TLS {
		config.Certificates[index], err = tls.LoadX509KeyPair(cert.Cert, cert.Key)
		if err != nil {
			return err
		}
	}

	config.BuildNameToCertificate()

	conn, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}

	listener := tls.NewListener(conn, config)
	return server.Serve(listener)
}
