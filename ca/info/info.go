package info

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type Cert struct {
	cert []byte
}

func NewCert(cert []byte) *Cert {
	return &Cert{cert: cert}
}

func (i *Cert) GetInfo() (*x509.Certificate, error) {
	block, _ := pem.Decode(i.cert)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("failed to decode PEM block containing certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	return cert, err
}

type Crl struct {
	crl []byte
}

func NewCrl(crl []byte) *Crl {
	return &Crl{crl: crl}
}

func (i *Crl) GetInfo() (*x509.RevocationList, error) {
	block, _ := pem.Decode(i.crl)
	if block == nil || block.Type != "X509 CRL" {
		return nil, errors.New("failed to decode PEM block containing CRL")
	}
	crl, err := x509.ParseRevocationList(block.Bytes)
	return crl, err
}
