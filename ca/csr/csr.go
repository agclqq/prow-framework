package csr

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"net"
)

type Csr struct {
	prvKey []byte
	SubjC  []string //country
	SubjST []string //province
	SubjL  []string //locality
	SubjO  []string //organization
	SubjOU []string //organizational unit
	SubjCN string   //common name
	Days   int
	Dns    []string
	Ips    []string
}
type CsrOption func(*Csr)

func NewCsr(prvKey []byte, SubjC, SubjST, SubjL, SubjO, SubjOU []string, SubjCN string, opts ...CsrOption) *Csr {
	c := &Csr{
		prvKey: prvKey,
		SubjC:  SubjC,
		SubjST: SubjST,
		SubjL:  SubjL,
		SubjO:  SubjO,
		SubjOU: SubjOU,
		SubjCN: SubjCN,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithDns(dns []string) CsrOption {
	return func(c *Csr) {
		c.Dns = dns
	}
}
func WithIps(ips []string) CsrOption {
	return func(c *Csr) {
		c.Ips = ips
	}
}

func (c *Csr) Gen() ([]byte, error) {
	//解析ca私钥
	block, _ := pem.Decode(c.prvKey)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	var prvKey *rsa.PrivateKey
	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		prvKey = key
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not an RSA private key")
		}
		prvKey = rsaKey
	default:
		return nil, errors.New("not an RSA private key")
	}

	csr := &x509.CertificateRequest{
		Subject: pkix.Name{
			Country:            c.SubjC,
			Organization:       c.SubjO,
			OrganizationalUnit: c.SubjOU,
			Locality:           c.SubjL,
			Province:           c.SubjST,
			CommonName:         c.SubjCN,
		},
	}
	if c.Dns != nil {
		csr.DNSNames = c.Dns
	}
	if c.Ips != nil {
		for _, ip := range c.Ips {
			addr := net.ParseIP(ip)
			if addr == nil {
				return nil, errors.New("invalid IP address")
			}
			csr.IPAddresses = append(csr.IPAddresses, addr)
		}
	}
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csr, prvKey)
	if err != nil {
		return nil, err
	}
	csrPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
	return csrPem, nil
}
