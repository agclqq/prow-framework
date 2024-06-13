package revoke

import (
	"errors"
	"fmt"
	"strings"

	"github.com/agclqq/prow-framework/execcmd"

	"rms/infra/ca"
)

type Revoke struct {
	caKeyPath  string
	caPemPath  string
	caConfPath string
	certPath   string
	crlPath    string
}

type Option func(*Revoke)

func NewRevoke(caKeyPath, caPemPath, caConfPath, certPath, crlPath string, opts ...Option) *Revoke {
	r := &Revoke{
		caKeyPath:  caKeyPath,
		caPemPath:  caPemPath,
		caConfPath: caConfPath,
		certPath:   certPath,
		crlPath:    crlPath,
	}
	for _, opt := range opts {
		opt(r)
	}

	return r
}
func (r *Revoke) verify() error {
	if !ca.SslVerify() {
		return errors.New("openssl is not installed")
	}
	if r.caConfPath == "" {
		return errors.New("ca conf path is empty")
	}
	if r.caKeyPath == "" {
		return errors.New("ca key path is empty")

	}
	if r.caPemPath == "" {
		return errors.New("ca pem path is empty")
	}
	if r.certPath == "" {
		return errors.New("the path of the certificate to be revoked cannot be empty")
	}
	return nil
}
func (r *Revoke) Revoke() error {
	err := r.verify()
	if err != nil {
		return err
	}
	// openssl ca -revoke certs/client.pem -keyfile ca.key -cert ca.pem -config ca.conf
	if log, err := execcmd.Command("openssl", "ca", "-revoke", r.certPath, "-keyfile", r.caKeyPath, "-cert", r.caPemPath, "-config", r.caConfPath); err != nil {
		fmt.Printf("Revoke error: %s\n", log)
		if strings.Contains(string(log), "ERROR:Already revoked") {
			return errors.New("the certificate has been revoked")
		}
		return err
	}

	if log, err := execcmd.Command("openssl", "ca", "-gencrl", "-keyfile", r.caKeyPath, "-cert", r.caPemPath, "-config", r.caConfPath); err != nil {
		fmt.Printf("Revoke error: %s\n", log)
		return err
	}
	return nil
}
