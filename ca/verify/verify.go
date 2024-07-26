package verify

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type Cert struct {
	caCert        []byte
	cert          []byte
	intermediates [][]byte
}

type CertOption func(*Cert)

func NewCert(caCert, cert []byte, opts ...CertOption) *Cert {
	c := &Cert{cert: cert}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithIntermediates(intermediates [][]byte) CertOption {
	return func(c *Cert) {
		c.intermediates = intermediates
	}
}

func (c *Cert) verify() error {
	caCert, err := c.parseCert(c.caCert)
	if err != nil {
		return err
	}
	certPool := x509.NewCertPool()
	certPool.AddCert(caCert)

	intermediatePool := x509.NewCertPool()
	for _, intermediate := range c.intermediates {
		cert, err := c.parseCert(intermediate)
		if err != nil {
			return err
		}
		intermediatePool.AddCert(cert)
	}

	cert, err := c.parseCert(c.cert)
	if err != nil {
		return err
	}

	_, err = cert.Verify(x509.VerifyOptions{
		Roots:         certPool,
		Intermediates: intermediatePool,
	})
	return err

}

func (c *Cert) parseCert(certData []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certData)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("failed to decode PEM block containing certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	return cert, err
}
