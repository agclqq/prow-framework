package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	prvKey []byte
	pubKey []byte
}

var (
	svrKey      []byte
	svrCert     []byte
	svrKeyPath  string
	svrCertPath string
	caCertByte  []byte
)

func ctlTest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello")) // #nosec G104
}

func router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /test", ctlTest)
	return mux
}

func Svr(caCert []byte, certPath, keyPath string) error {
	if caCert == nil {
		return errors.New("ca cert is nil")
	}
	caCertByte = caCert
	err := CreateSvrCert(caCert, certPath, keyPath, "svr")
	if err != nil {
		return err
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return errors.New("failed to parse CA certificate")
	}
	pair, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr:              ":8081",
		Handler:           router(),
		IdleTimeout:       75 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		TLSConfig: &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{pair},
			RootCAs:      certPool,
			//ClientAuth:         0,
			//InsecureSkipVerify: false,
		},
	}
	go func() {
		fmt.Printf("start svr https server at: %s\n", server.Addr)
		if err := server.ListenAndServeTLS("", ""); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				fmt.Println(err)
				return
			}
			_ = fmt.Errorf("start ca http server is error: %s\n", err)
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	server.Shutdown(context.Background()) // #nosec G104
	return nil
}

func (s *Server) Crt() []byte {
	return nil
}
