package rabbit

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
)

func (a *rabbit) Listen() (err error) {
	select {
	case err = <-a.ListenHTTP():
		return
	case err = <-a.ListenHTTPS():
		return
	}
}

func (a *rabbit) ListenHTTP() <-chan error {
	e := make(chan error)
	go func() {
		port := fmt.Sprintf(":%v", a.config.Ports.HTTP)
		a.Log("listening http on %v", port)
		e <- http.ListenAndServe(port, a)
	}()
	return e
}

func (a *rabbit) ListenHTTPS() <-chan error {
	e := make(chan error)
	if a.config.Ports.HTTPS != 0 && len(a.config.Certs) > 0 {
		go func() {
			port := fmt.Sprintf(":%v", a.config.Ports.HTTPS)
			a.Log("listening https on %v", port)
			e <- a.listenHTTPS(port)
		}()
	} else {
		a.Log("https/certs not configured")
	}
	return e
}

func (a *rabbit) listenHTTPS(port string) error {
	config := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		NextProtos:               []string{"h2", "http/1.1"},
	}

	for _, cert := range a.config.Certs {
		a.Log("loading tls certificate (%v, %v)", cert.Cert, cert.Key)
		kp, err := tls.LoadX509KeyPair(cert.Cert, cert.Key)
		if err != nil {
			return err
		}
		config.Certificates = append(config.Certificates, kp)
	}

	config.BuildNameToCertificate()

	server := &http.Server{
		Addr:      port,
		Handler:   a,
		TLSConfig: config,
	}

	conn, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}

	listener := tls.NewListener(conn, config)
	return server.Serve(listener)
}
