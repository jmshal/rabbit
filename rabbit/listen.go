package rabbit

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"strconv"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

func getPortsAndHandler(a *Rabbit) (http string, https string, handler http.Handler) {
	http = ":" + strconv.Itoa(int(a.config.Ports.HTTP))
	https = ":" + strconv.Itoa(int(a.config.Ports.HTTPS))
	handler = bugsnag.Handler(a.mux)
	return
}

func (a *Rabbit) listenHTTP() error {
	port, _, handler := getPortsAndHandler(a)
	return http.ListenAndServe(port, handler)
}

func (a *Rabbit) listenTLS() error {
	_, port, handler := getPortsAndHandler(a)

	config := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		NextProtos:               []string{"h2", "h2-14", "http/1.1"},
	}

	var err error
	config.Certificates = make([]tls.Certificate, len(a.config.Certs))
	for index, cert := range a.config.Certs {
		config.Certificates[index], err = tls.LoadX509KeyPair(cert.Cert, cert.Key)
		if err != nil {
			return err
		}
	}

	config.BuildNameToCertificate()

	server := &http.Server{
		Addr:      port,
		Handler:   handler,
		TLSConfig: config,
	}

	conn, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}

	listener := tls.NewListener(conn, config)
	return server.Serve(listener)
}

func (a *Rabbit) Listen() error {
	errs := make(chan error)
	if a.config.Ports.HTTPS != 0 {
		go func() {
			log.Printf("Listening tls server on :%v", a.config.Ports.HTTPS)
			errs <- a.listenTLS()
		}()
	}
	go func() {
		log.Printf("Listening http server on :%v", a.config.Ports.HTTP)
		errs <- a.listenHTTP()
	}()
	return <-errs
}
