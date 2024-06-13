package sign

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/agclqq/prow-framework/execcmd"

	"rms/infra/ca"
)

type Cert struct {
	caKey  string
	caPem  string
	Dir    string
	csr    []byte
	prvKey string
	SubjC  string
	SubjST string
	SubjL  string
	SubjO  string
	SubjOU string
	SubjCN string
	Days   int
	Dns    string
	Domain string
	CaType string
}
type CertOption func(*Cert)

func NewCert(caKeyPath, caPemPath string, opts ...CertOption) *Cert {
	c := &Cert{
		caKey:  caKeyPath,
		caPem:  caPemPath,
		Days:   365,
		CaType: "client",
	}
	for _, opt := range opts {
		opt(c)
	}

	return c
}
func WithCountry(country string) CertOption {
	return func(c *Cert) {
		c.SubjC = country
	}
}

func WithState(state string) CertOption {
	return func(c *Cert) {
		c.SubjST = state
	}
}

func WithLocality(locality string) CertOption {
	return func(c *Cert) {
		c.SubjL = locality
	}
}

func WithOrganization(organization string) CertOption {
	return func(c *Cert) {
		c.SubjO = organization
	}
}

func WithOrganizationalUnit(organizationalUnit string) CertOption {
	return func(c *Cert) {
		c.SubjOU = organizationalUnit
	}
}

func WithCommonName(commonName string) CertOption {
	return func(c *Cert) {
		c.SubjCN = commonName
	}
}

func WithDir(dir string) CertOption {
	return func(c *Cert) {
		c.Dir = dir
	}
}
func WithPrvKey(key string) CertOption {
	return func(c *Cert) {
		c.prvKey = key
	}
}
func WithDays(days int) CertOption {
	return func(c *Cert) {
		c.Days = days
	}
}
func WithDns(dns string) CertOption {
	return func(c *Cert) {
		c.Dns = dns
	}
}

func WithServerCaType() CertOption {
	return func(c *Cert) {
		c.CaType = "server"
	}
}
func WithCsr(csr []byte) CertOption {
	return func(c *Cert) {
		c.csr = csr
	}
}

func (c *Cert) createConf() error {
	if !ca.SslVerify() {
		return errors.New("openssl is not installed")
	}
	if c.CaType == "server" {
		if c.Dns == "" {
			return errors.New("dns cannot be empty")
		}
	}

	var err error
	if c.Dir == "" {
		c.Dir, err = os.MkdirTemp("", "ca_"+c.CaType)
		if err != nil {
			fmt.Println("Failed to create temporary directory")
			return err
		}
	} else {
		err = os.MkdirAll(c.Dir, 0755)
		if err != nil {
			fmt.Println("Failed to create directory:", c.Dir)
			return err
		}
	}

	confFilePath := filepath.Join(c.Dir, c.CaType+".conf")
	var confContent string
	if c.CaType == "server" {
		confContent = fmt.Sprintf(ca.SignServerCertConf(), c.SubjC, c.SubjST, c.SubjL, c.SubjO, c.SubjOU, c.SubjCN, c.Dns)
	} else {
		confContent = fmt.Sprintf(ca.SignClientCertConf(), c.SubjC, c.SubjST, c.SubjL, c.SubjO, c.SubjOU, c.SubjCN)
	}
	err = os.WriteFile(confFilePath, []byte(confContent), 0644)
	return err
}

// Sign signs the certificate, returns the path of the key and the certificate
func (c *Cert) Sign() (string, string, error) {
	err := c.createConf()
	if err != nil {
		return "", "", err
	}
	keyPath := filepath.Join(c.Dir, c.CaType+".key")
	csrPath := filepath.Join(c.Dir, c.CaType+".csr")
	confFilePath := filepath.Join(c.Dir, c.CaType+".conf")
	pemPath := filepath.Join(c.Dir, c.CaType+".pem")

	if c.caPem == "" || c.caKey == "" {
		return "", "", errors.New("ca pem or key path cannot be empty")
	}

	var log []byte
	if c.prvKey == "" {
		// openssl genrsa -out server.key 2048
		if log, err = execcmd.Command("openssl", "genrsa", "-out", keyPath, "2048"); err != nil {
			fmt.Printf("Failed to generate RSA key:%v,log:%s", err, log)
			return "", "", err
		}
	} else {
		err = os.WriteFile(keyPath, []byte(c.prvKey), 0644)
		if err != nil {
			fmt.Println("Failed to write private key:", err)
			return "", "", err
		}
	}

	// openssl req -new -sha256 -out server.csr -key server.key -config server.conf
	if log, err = execcmd.Command("openssl", "req", "-new", "-sha256", "-out", csrPath, "-key", keyPath, "-config", confFilePath); err != nil {
		fmt.Printf("Failed to generate CSR:%v,log:%s", err, log)
		return "", "", err
	}

	if c.CaType == "server" {
		// openssl x509 -req -days 365 -CA ca.pem -CAkey private/ca.key -CAcreateserial -in server.csr -out server.pem -extensions req_ext -extfile server.conf
		log, err = execcmd.Command("openssl", "x509", "-req", "-days", fmt.Sprintf("%d", c.Days), "-CA", c.caPem, "-CAkey", c.caKey, "-CAcreateserial", "-in", csrPath, "-out", pemPath, "-extensions", "req_ext", "-extfile", confFilePath)
	} else {
		// openssl x509 -req -days 365 -CA ca.pem -CAkey private/ca.key -CAcreateserial -in client.csr -out client.pem
		log, err = execcmd.Command("openssl", "x509", "-req", "-days", fmt.Sprintf("%d", c.Days), "-CA", c.caPem, "-CAkey", c.caKey, "-CAcreateserial", "-in", csrPath, "-out", pemPath)
	}

	if err != nil {
		fmt.Printf("Failed to sign certificate:%v,log:%s", err, log)
	}
	return keyPath, pemPath, err
}

func (c *Cert) signFromCsr() (string, string, error) {
	err := c.createConf()
	if err != nil {
		return "", "", err
	}

	if c.csr == nil {
		return "", "", errors.New("csr cannot be empty")
	}
	if c.CaType == "server" && c.Dns == "" {
		return "", "", errors.New("dns cannot be empty")
	}
	if c.CaType == "" {
		c.CaType = "client"
	}
	err = os.WriteFile(filepath.Join(c.Dir, c.CaType+".csr"), c.csr, 0644)
	if err != nil {
		return "", "", err
	}

	csrPath := filepath.Join(c.Dir, c.CaType+".csr")
	confFilePath := filepath.Join(c.Dir, c.CaType+".conf")
	pemPath := filepath.Join(c.Dir, c.CaType+".pem")

	var log []byte
	if c.CaType == "server" {
		// openssl x509 -req -days 365 -CA ca.pem -CAkey private/ca.key -CAcreateserial -in server.csr -out server.pem -extensions req_ext -extfile server.conf
		log, err = execcmd.Command("openssl", "x509", "-req", "-days", fmt.Sprintf("%d", c.Days), "-CA", c.caPem, "-CAkey", c.caKey, "-CAcreateserial", "-in", csrPath, "-out", pemPath, "-extensions", "req_ext", "-extfile", confFilePath)
	} else {
		// openssl x509 -req -days 365 -CA ca.pem -CAkey private/ca.key -CAcreateserial -in client.csr -out client.pem
		log, err = execcmd.Command("openssl", "x509", "-req", "-days", fmt.Sprintf("%d", c.Days), "-CA", c.caPem, "-CAkey", c.caKey, "-CAcreateserial", "-in", csrPath, "-out", pemPath)
	}
	if err != nil {
		fmt.Printf("Failed to sign certificate:%v,log:%s", err, log)
	}
	return "", pemPath, err
}
