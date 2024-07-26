package ca

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/agclqq/prow-framework/ca/issuance"
)

var (
	caKey      []byte
	caCert     []byte
	caCertPath string
)

type postData struct {
	Csr []byte `json:"csr"`
}
type response struct {
	status int
}

func ctlReqCert(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	pd := &postData{}
	//fmt.Println(string(body))
	err = json.Unmarshal(body, pd)
	//fmt.Printf("csr: %s\n", pd.Csr)
	if err != nil {
		//fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	issueType := issuance.IssueTypeClient
	t := request.URL.Query().Get("t")
	if t == "1" {
		issueType = issuance.IssueTypeServer
	}
	cert, err := issuance.NewCert(caKey, caCert, pd.Csr, issuance.WithIssueType(issueType)).Sign()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	//resp := response{status: http.StatusOK}
	//jsonResp, err := json.Marshal(resp)
	//if err != nil {
	//	writer.WriteHeader(http.StatusBadRequest)
	//	writer.Write([]byte(err.Error()))
	//	return
	//}
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "cert"))
	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(cert)))
	writer.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	writer.Write(cert)
}
func ctlCert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "cert"))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(caCert)))
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Write(caCert)
}
func router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /cert", ctlCert)
	mux.HandleFunc("POST /reqCert", ctlReqCert)
	return mux
}
func Svr(caSvrCertPath, caSvrKeyPath string) error {
	err := CreateCaCert("ca.crt", "ca.key")
	if err != nil {
		return err
	}

	caSvrCert, caSvrKey, err := IssueCaServerCert()
	if err != nil {
		return err
	}
	err = os.WriteFile(caSvrCertPath, caSvrCert, 0666)
	if err != nil {
		return err
	}
	err = os.WriteFile(caSvrKeyPath, caSvrKey, 0600)
	if err != nil {
		return err
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return errors.New("failed to parse CA certificate")
	}
	pair, err := tls.LoadX509KeyPair(caSvrCertPath, caSvrKeyPath)
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr:        ":8080",
		Handler:     router(),
		IdleTimeout: 75 * time.Second,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{pair},
			RootCAs:      certPool,
			//ClientAuth:         0,
			//InsecureSkipVerify: false,
		},
	}
	go func() {
		fmt.Printf("start ca https server at: %s\n", server.Addr)
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
	server.Shutdown(context.Background())
	return nil
}
