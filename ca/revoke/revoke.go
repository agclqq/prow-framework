package revoke

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"time"
)

type Revoke struct {
	caCert []byte
	caKey  []byte
	cert   []byte
	crl    []byte
}

type RevokeOption func(*Revoke)

func NewRevoke(caCert, caKey, cert []byte, opts ...RevokeOption) *Revoke {
	r := &Revoke{
		caCert: caCert,
		caKey:  caKey,
		cert:   cert,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func WithCrl(crl []byte) RevokeOption {
	return func(r *Revoke) {
		r.crl = crl
	}
}

func (r *Revoke) Revoke() ([]byte, error) {
	//解析ca证书
	block, _ := pem.Decode(r.caCert)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing certificate")
	}
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	//解析ca私钥
	block, _ = pem.Decode(r.caKey)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	caKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//解析待吊销证书
	block, _ = pem.Decode(r.cert)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	rl := &x509.RevocationList{}
	if r.crl != nil {
		block, _ = pem.Decode(r.crl)
		if block == nil {
			return nil, errors.New("failed to decode PEM block containing CRL")
		}
		rl, err = x509.ParseRevocationList(block.Bytes)
		if err != nil {
			return nil, err
		}

	}
	rl.RevokedCertificateEntries = append(rl.RevokedCertificateEntries, x509.RevocationListEntry{
		SerialNumber:   cert.SerialNumber,
		RevocationTime: time.Now(),
	})
	if rl.Number == nil {
		rl.Number = big.NewInt(1)
	} else {
		rl.Number = big.NewInt(0).Add(rl.Number, big.NewInt(1))
	}

	rl.ThisUpdate = time.Now()
	rl.NextUpdate = time.Now().AddDate(0, 0, 7)
	crlBytes, err := x509.CreateRevocationList(rand.Reader, rl, caCert, caKey)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "X509 CRL", Bytes: crlBytes}), nil
}
