package selfsign

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/agclqq/prow-framework/execcmd"
	"github.com/agclqq/prow-framework/file"

	"rms/infra/ca"
)

type Ca struct {
	dir    string
	SubjC  string
	SubjST string
	SubjL  string
	SubjO  string
	SubjOU string
	SubjCN string
	Days   int
}
type CaOption func(*Ca)

func NewCa(c, st, l, o, ou, cn string, opts ...CaOption) *Ca {
	ca := &Ca{
		SubjC:  c,
		SubjST: st,
		SubjL:  l,
		SubjO:  o,
		SubjOU: ou,
		SubjCN: cn,
		Days:   365,
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
func WithDir(dir string) CaOption {
	return func(c *Ca) {
		c.dir = dir
	}
}
func (c *Ca) createConf() error {
	if !ca.SslVerify() {
		return errors.New("openssl is not installed")
	}
	if c.dir == "" {
		baseDir := filepath.Join(os.TempDir(), "ca")
		err := os.MkdirAll(baseDir, 0755)
		if err != nil {
			return err
		}
		dir, err := os.MkdirTemp(baseDir, "ca-")
		if err != nil {
			return err
		}
		c.dir = dir
	}
	if !file.Exist(c.dir) {
		err := os.MkdirAll(c.dir, 0755)
		if err != nil {
			return err
		}
	}

	confFilePath := filepath.Join(c.dir, "ca.conf")
	fileContent := fmt.Sprintf(ca.SelfSignCaConf(), c.dir, c.SubjC, c.SubjST, c.SubjL, c.SubjO, c.SubjOU, c.SubjCN)
	err := os.WriteFile(confFilePath, []byte(fileContent), 0644)
	return err
}

// Sign signs the certificate
// return keyPath, pemPath, error
func (c *Ca) Sign() (string, string, error) {
	err := c.createConf()
	if err != nil {
		return "", "", err
	}
	keyPath := filepath.Join(c.dir, "ca.key")
	csrPath := filepath.Join(c.dir, "ca.csr")
	confPath := filepath.Join(c.dir, "ca.conf")
	pemPath := filepath.Join(c.dir, "ca.pem")

	// openssl genrsa -out server.key 2048
	if log, err := execcmd.Command("openssl", "genrsa", "-out", keyPath, "2048"); err != nil {
		fmt.Printf("Failed to generate RSA key:%v,%s \n", err, log)
		return "", "", err
	}
	// openssl req -new -sha256 -out ca.csr -key ca.key -config ca.conf
	if log, err := execcmd.Command("openssl", "req", "-new", "-sha256", "-out", csrPath, "-key", keyPath, "-config", confPath); err != nil {
		fmt.Printf("Failed to generate CSR:%v,log:%s \n", err, log)
		return "", "", err
	}

	if log, err := execcmd.Command("openssl", "x509", "-req", "-days", fmt.Sprintf("%d", c.Days), "-in", csrPath, "-signkey", keyPath, "-out", pemPath); err != nil {
		fmt.Printf("Failed to generate PEM:%v,log:%s", err, log)
		return "", "", err
	}
	return keyPath, pemPath, nil
}
