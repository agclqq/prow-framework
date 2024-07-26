package server

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/agclqq/prow-framework/ca/csr"
	"github.com/agclqq/prow-framework/ca/prvkey"
	"github.com/agclqq/prow-framework/file"
)

func CreateSvrCert(caCert []byte, certPath, keyPath, cn string) error {
	if file.Exist(keyPath) && file.Exist(certPath) {
		return nil
	}
	//创建服务端私钥
	svrKey, err := prvkey.Gen(2048)
	if err != nil {
		return err
	}
	//创建请求证书
	csrByte, err := csr.NewCsr(svrKey, []string{"CN"}, []string{"BJ"}, []string{"BJ"}, []string{"company"}, []string{"a group"}, cn, csr.WithIps([]string{"127.0.0.1"})).Gen()
	if err != nil {
		return err
	}
	//申请证书
	svrCert, err := reqCert(caCert, csrByte)
	if err != nil {
		return err
	}

	kpd := path.Dir(keyPath)
	if !file.Exist(kpd) {
		err = os.MkdirAll(kpd, 0666)
		if err != nil {
			return err
		}
	}
	err = os.WriteFile(keyPath, svrKey, 0660)
	if err != nil {
		return err
	}

	cpd := path.Dir(certPath)
	if !file.Exist(cpd) {
		err = os.MkdirAll(cpd, 0666)
		if err != nil {
			return err
		}
	}
	err = os.WriteFile(certPath, svrCert, 0666)
	if err != nil {
		return err
	}
	return nil
}

func reqCert(caCert, csr []byte) ([]byte, error) {
	if len(csr) == 0 {
		return nil, errors.New("csr is nil")
	}
	type csrStruct struct {
		Csr []byte `json:"csr"`
	}

	c := &csrStruct{
		Csr: csr,
	}

	csrJson, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("failed to parse CA certificate")
	}
	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}
	resp, err := cli.Post("https://127.0.0.1:8080/reqCert?t=1", "application/json", bytes.NewBuffer(csrJson))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, err
	}

	cert, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return cert, nil
}
