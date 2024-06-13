package info

import (
	"errors"
	"os"
	"strings"

	"github.com/agclqq/prow-framework/execcmd"
	"github.com/agclqq/prow-framework/file"
)

type Cert struct {
	CertPath    string
	CertContent string
	tmpDir      string
}

type CertOption func(*Cert)

func NewCert(opts ...CertOption) *Cert {
	c := &Cert{}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithCertPath(certPath string) CertOption {
	return func(c *Cert) {
		c.CertPath = certPath
	}
}
func WithCertContent(certContent string) CertOption {
	return func(c *Cert) {
		c.CertContent = certContent
	}
}

func (c *Cert) verify() error {
	if c.CertPath == "" && c.CertContent == "" {
		return errors.New("证书路径和证书内容不能同时为空")
	}
	if c.CertPath == "" {
		temp, err := os.MkdirTemp("", "cert-*")
		if err != nil {
			return err
		}
		c.tmpDir = temp
		err = os.WriteFile(temp+"/cert.pem", []byte(c.CertContent), 0644)
		if err != nil {
			return err
		}
		c.CertPath = temp + "/cert.pem"
	} else {
		if !file.Exist(c.CertPath) {
			return errors.New("证书路径不存在")
		}
	}
	return nil
}

func (c *Cert) GetSerial() (string, error) {
	if err := c.verify(); err != nil {
		return "", err
	}
	defer os.RemoveAll(c.tmpDir)

	log, err := execcmd.Command("openssl", "x509", "-in", c.CertPath, "-noout", "-serial")
	if err != nil {
		return "", err
	}
	rs := strings.Split(string(log), "=")
	if len(rs) != 2 {
		return "", errors.New("获取序列号失败")
	}
	return strings.TrimSpace(rs[1]), nil
}

func (c *Cert) GetStartTime() (string, error) {
	if err := c.verify(); err != nil {
		return "", err
	}
	defer os.RemoveAll(c.tmpDir)

	log, err := execcmd.Command("openssl", "x509", "-in", c.CertPath, "-noout", "-startdate")
	if err != nil {
		return "", err
	}
	rs := strings.Split(string(log), "=")
	if len(rs) != 2 {
		return "", errors.New("获取证书开始时间失败")
	}
	return strings.TrimSpace(rs[1]), nil
}

func (c *Cert) GetEndTime() (string, error) {
	if err := c.verify(); err != nil {
		return "", err
	}
	defer os.RemoveAll(c.tmpDir)

	log, err := execcmd.Command("openssl", "x509", "-in", c.CertPath, "-noout", "-enddate")
	if err != nil {
		return "", err
	}
	rs := strings.Split(string(log), "=")
	if len(rs) != 2 {
		return "", errors.New("获取证书结束时间失败")
	}
	return strings.TrimSpace(rs[1]), nil
}
