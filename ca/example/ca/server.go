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
	"path/filepath"
	"time"

	"github.com/agclqq/prow-framework/ca/info"
	"github.com/agclqq/prow-framework/ca/issuance"
	"github.com/agclqq/prow-framework/ca/ocsp"
	"github.com/agclqq/prow-framework/ca/revoke"

	ocsp1 "golang.org/x/crypto/ocsp"
)

var (
	caKey      []byte
	caCert     []byte
	caCertPath string = "ca.crt"
	caKeyPath  string = "ca.key"
	ctlPath    string = "ctl.pem"
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
		writer.Write([]byte(err.Error())) // #nosec G104
		return
	}
	pd := &postData{}
	//fmt.Println(string(body))
	err = json.Unmarshal(body, pd)
	//fmt.Printf("csr: %s\n", pd.Csr)
	if err != nil {
		//fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error())) // #nosec G104
		return
	}
	issueType := issuance.IssueTypeClient
	t := request.URL.Query().Get("t")
	if t == "1" {
		issueType = issuance.IssueTypeServer
	}
	cert, err := issuance.NewCert(caCert, caKey, pd.Csr, issuance.WithIssueType(issueType), issuance.WithOcspServer([]string{"https://127.0.0.1:8080/ocsp"}), issuance.WithCrlPoint([]string{"https://127.0.0.1:8080/crl"})).Sign()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error())) // #nosec G104
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "cert"))
	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(cert)))
	writer.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	writer.Write(cert) // #nosec G104
}
func ctlCert(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "cert"))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(caCert)))
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Write(caCert) // #nosec G104
}
func ctlOcsp(w http.ResponseWriter, r *http.Request) {
	certPem, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to read cert,the certificate to be queried cannot be empty")) // #nosec G104
		return
	}
	x509Cert, err := info.NewCert(certPem).GetInfo()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) // #nosec G104
		return
	}
	x509CaCert, err := info.NewCert(caCert).GetInfo()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) // #nosec G104
		return
	}
	ctl, err := os.ReadFile(ctlPath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) // #nosec G104
		return
	}
	reqOcsp, err := ocsp1.CreateRequest(x509Cert, x509CaCert, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) // #nosec G104
		return
	}
	bytes, err := ocsp.Ocsp(reqOcsp, caCert, caKey, ctl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) // #nosec G104
		return
	}

	w.WriteHeader(http.StatusOK)
	//w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "ocsp"))
	//w.Header().Set("Content-Type", "application/octet-stream")
	//w.Header().Set("Content-Length", fmt.Sprintf("%d", len(bytes)))
	w.Write(bytes) // #nosec G104
}
func ctlCrl(w http.ResponseWriter, r *http.Request) {
	ctl, err := os.ReadFile(ctlPath)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "ctl.pem"))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(ctl)))
	w.Write(ctl) // #nosec G104
}
func ctlRevoke(w http.ResponseWriter, r *http.Request) {
	//TODO:身份验证必须要做，此处省略

	cert, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("failed to read cert,the certificate to be queried cannot be empty")) // #nosec G104
		return
	}
	crl, err := os.ReadFile(ctlPath)
	if err != nil {
		w.Write([]byte(err.Error())) // #nosec G104
		return
	}
	newCrl, err := revoke.NewRevoke(caCert, caKey, cert, revoke.WithCrl(crl)).Revoke()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) // #nosec G104
		return
	}
	err = os.WriteFile(ctlPath, newCrl, 0600)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// #nosec G104
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
func router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /cert", ctlCert)
	mux.HandleFunc("POST /reqCert", ctlReqCert)
	mux.HandleFunc("POST /revoke", ctlRevoke)
	mux.HandleFunc("POST /ocsp", ctlOcsp)
	mux.HandleFunc("GET /crl", ctlCrl)
	return mux
}
func Svr(caSvrCertPath, caSvrKeyPath string) error {
	caSvrCertPath = filepath.Clean(caSvrCertPath)
	caSvrKeyPath = filepath.Clean(caSvrKeyPath)
	err := CreateCaCert(caCertPath, caKeyPath)
	if err != nil {
		return err
	}

	caSvrCert, caSvrKey, err := IssueCaServerCert()
	if err != nil {
		return err
	}
	err = os.WriteFile(caSvrCertPath, caSvrCert, 0600)
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
		Addr:              ":8080",
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
	server.Shutdown(context.Background()) // #nosec G104
	return nil
}
