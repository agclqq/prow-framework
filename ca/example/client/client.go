package client

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

var (
	cliKey      []byte
	cliCert     []byte
	caCertByte  []byte
	cliKeyPath  string
	cliCertPath string
)

func Cli(caCert []byte, keyPath, certPath string) error {
	if caCert == nil {
		return errors.New("ca cert is nil")
	}
	caCertByte = caCert
	err := createCliCert(keyPath, certPath)
	if err != nil {
		return err
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return errors.New("failed to parse CA certificate")
	}
	pair, err := tls.X509KeyPair(cliCert, cliKey)
	if err != nil {
		return err
	}
	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{pair}, // 设置客户端证书和私钥
				RootCAs:      certPool,                // 设置服务端CA证书池用于验证服务端证书
				//ClientAuth:            tls.RequireAndVerifyClientCert, // 启用客户端证书验证（即双向证书验证）
				//VerifyPeerCertificate: Vpc,
			},
		},
	}
	_, err = cli.Get("https://127.0.0.1:8081/test")
	return err
}

func createCliCert(keyPath, certPath string) error {
	if file.Exist(keyPath) && file.Exist(certPath) {
		if len(cliKey) == 0 || len(cliCert) == 0 {
			key, err := os.ReadFile(keyPath)
			if err != nil {
				return err
			}
			cliKey = key

			cert, err := os.ReadFile(certPath)
			if err != nil {
				return err
			}
			cliCert = cert
		}
		return nil
	}
	//创建服务端私钥
	key, err2 := prvkey.Gen(2048)
	if err2 != nil {
		return err2
	}
	cliKey = key
	//创建请求证书
	csrByte, err := csr.NewCsr(cliKey, []string{"CN"}, []string{"BJ"}, []string{"BJ"}, []string{"company"}, []string{"a group"}, "127.0.0.1").Gen()
	if err != nil {
		return err
	}
	//申请证书
	cert, err := reqCert(csrByte)
	if err != nil {
		return err
	}
	cliCert = cert

	kpd := path.Dir(keyPath)
	if !file.Exist(kpd) {
		err = os.MkdirAll(kpd, 0666)
		if err != nil {
			return err
		}
	}
	err = os.WriteFile(keyPath, cliKey, 0660)
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
	err = os.WriteFile(certPath, cliCert, 0666)
	if err != nil {
		return err
	}
	return nil
}

func reqCert(csr []byte) ([]byte, error) {
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
	if ok := certPool.AppendCertsFromPEM(caCertByte); !ok {
		return nil, errors.New("failed to parse CA certificate")
	}
	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}
	resp, err := cli.Post("https://127.0.0.1:8080/reqCert", "application/json", bytes.NewBuffer(csrJson))
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
