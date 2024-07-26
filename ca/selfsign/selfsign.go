package selfsign

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

type Ca struct {
	SubjC  []string //country
	SubjST []string //province
	SubjL  []string //locality
	SubjO  []string //organization
	SubjOU []string //organizational unit
	SubjCN string   //common name
	Days   int
	keyBit int
}
type CaOption func(*Ca)

func NewCa(c, st, l, o, ou []string, cn string, opts ...CaOption) *Ca {
	ca := &Ca{
		SubjC:  c,
		SubjST: st,
		SubjL:  l,
		SubjO:  o,
		SubjOU: ou,
		SubjCN: cn,
		Days:   365,
		keyBit: 2048,
	}

	for _, opt := range opts {
		opt(ca)
	}
	return ca
}
func WithDays(days int) CaOption {
	return func(c *Ca) {
		c.Days = days
	}
}

func WithBit(bit int) CaOption {
	return func(c *Ca) {
		c.keyBit = bit
	}
}
func (c *Ca) Sign() ([]byte, []byte, error) {
	//创建私钥
	prvKey, err := rsa.GenerateKey(rand.Reader, c.keyBit)
	if err != nil {
		return nil, nil, err
	}
	prvKeyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(prvKey)})
	subj := pkix.Name{
		Country:            c.SubjC,
		Organization:       c.SubjO,
		OrganizationalUnit: c.SubjOU,
		Locality:           c.SubjL,
		Province:           c.SubjST,
		StreetAddress:      nil,
		PostalCode:         nil,
		SerialNumber:       "",
		CommonName:         c.SubjCN,
		Names:              nil,
		ExtraNames:         nil,
	}

	//创建CSR
	csrTpl := x509.CertificateRequest{
		SignatureAlgorithm: x509.SHA256WithRSA,
		Subject:            subj,
	}
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &csrTpl, prvKey)
	if err != nil {
		return nil, nil, err
	}
	//csrPem:=pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
	//创建证书

	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		return nil, nil, err
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}
	certTpl := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               csr.Subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, c.Days),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	cert, err := x509.CreateCertificate(rand.Reader, &certTpl, &certTpl, prvKey.Public(), prvKey)
	if err != nil {
		return nil, nil, err
	}
	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert})
	return certPem, prvKeyPem, nil
}
