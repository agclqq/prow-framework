package info

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

type Cert struct {
	cert []byte
}
type CertInfos struct {
	*x509.Certificate //继承x509.Certificate
	Sha256            string
}

func NewCert(cert []byte) *Cert {
	return &Cert{cert: cert}
}

func (i *Cert) Get() (*CertInfos, error) {
	block, _ := pem.Decode(i.cert)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("failed to decode PEM block containing certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	digest := sha256.Sum256(cert.Raw)
	return &CertInfos{Certificate: cert, Sha256: hex.EncodeToString(digest[:])}, nil
}

type Crl struct {
	crl []byte
}

func NewCrl(crl []byte) *Crl {
	return &Crl{crl: crl}
}

func (i *Crl) Get() (*x509.RevocationList, error) {
	block, _ := pem.Decode(i.crl)
	if block == nil || block.Type != "X509 CRL" {
		return nil, errors.New("failed to decode PEM block containing CRL")
	}
	crl, err := x509.ParseRevocationList(block.Bytes)
	if err != nil {
		return nil, err
	}

	return crl, nil
}
