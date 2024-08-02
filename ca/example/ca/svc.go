package ca

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path"

	"github.com/agclqq/prow-framework/ca/crl"
	"github.com/agclqq/prow-framework/ca/csr"
	"github.com/agclqq/prow-framework/ca/issuance"
	"github.com/agclqq/prow-framework/ca/prvkey"
	"github.com/agclqq/prow-framework/ca/selfsign"
	"github.com/agclqq/prow-framework/file"
)

// CreateCaCert 创建CA证书
func CreateCaCert(certPath, keyPath string) error {
	if file.Exist(keyPath) && file.Exist(certPath) {
		caCertPath = certPath
		if len(caKey) == 0 || len(caCert) == 0 {
			key, err := os.ReadFile(keyPath)
			if err != nil {
				return err
			}
			caKey = key

			cert, err := os.ReadFile(certPath)
			if err != nil {
				return err
			}
			caCert = cert
		}
		return nil
	}
	//自签
	cert, key, err := selfsign.NewCa([]string{"CN"}, []string{"bj"}, []string{"bj"}, []string{""}, []string{""}, "my_company").Sign()
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
	cpd := path.Dir(certPath)
	if !file.Exist(cpd) {
		err = os.MkdirAll(cpd, 0666)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(keyPath, key, 0600)
	if err != nil {
		return err
	}
	err = os.WriteFile(certPath, cert, 0664)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(cert)
	if block == nil || block.Type != "CERTIFICATE" {
		return errors.New("failed to decode PEM block containing certificate")
	}
	caCertPem, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	block, _ = pem.Decode(key)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return errors.New("failed to decode PEM block containing private key")
	}
	caKeyPem, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	ctlPem, err := crl.InitCrl(caCertPem, caKeyPem)
	if err != nil {
		return err
	}
	err = os.WriteFile(ctlPath, ctlPem, 0664)
	if err != nil {
		return err
	}

	caKey = key
	caCert = cert

	return nil
}

// IssueCaServerCert 颁发ca服务器证书
// return prvKey,cert,error
func IssueCaServerCert() ([]byte, []byte, error) {
	key, err := prvkey.Gen(2048)
	if err != nil {
		return nil, nil, err
	}
	csrByte, err := csr.NewCsr(key, []string{"CN"}, []string{"CN"}, []string{"CN"}, []string{"CN"}, []string{"CN"}, "", csr.WithIps([]string{"127.0.0.1"})).Gen()
	if err != nil {
		return nil, nil, err
	}
	cert, err := issuance.NewCert(caCert, caKey, csrByte, issuance.WithIssueType(issuance.IssueTypeServer)).Sign()
	return cert, key, err
}

func IssueCert(caKeyPath, caCertPath string, csr []byte) ([]byte, error) {

	return issuance.NewCert(caCert, caKey, csr).Sign()
}
