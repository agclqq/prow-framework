package ocsp

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ocsp"
)

func Ocsp(body, caCert, caPrvKey, crl []byte) ([]byte, error) {
	req, err := ocsp.ParseRequest(body)
	if err != nil {
		return nil, err
	}
	// 解析CA证书
	caCertBlock, _ := pem.Decode(caCert)
	if caCertBlock == nil {
		return nil, errors.New("failed to decode PEM block containing certificate")
	}
	caCertPem, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return nil, err
	}
	// 解析CA私钥
	caPrvKeyBlock, _ := pem.Decode(caPrvKey)
	if caPrvKeyBlock == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	caPrvKeyPem, err := x509.ParsePKCS1PrivateKey(caPrvKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	// 解析吊销列表
	ctlBlock, _ := pem.Decode(crl)
	if ctlBlock == nil {
		return nil, errors.New("failed to decode PEM block containing CRL")
	}
	list, err := x509.ParseRevocationList(ctlBlock.Bytes)
	if err != nil {
		return nil, err
	}
	// 生成OCSP响应
	resp, err := generateOCSPResponse(req, caCertPem, caPrvKeyPem, list)
	return resp, err
}

// 生成OCSP响应
func generateOCSPResponse(req *ocsp.Request, caCert *x509.Certificate, caPrvKey crypto.Signer, crl *x509.RevocationList) ([]byte, error) {
	// 生成OCSP响应
	ocspResponseTemplate := ocsp.Response{
		Status:       ocsp.Good,
		SerialNumber: req.SerialNumber,
		ThisUpdate:   time.Now(),
		NextUpdate:   time.Now().AddDate(0, 0, 7),
		Certificate:  caCert,
	}
	// 检查请求的证书是否被吊销

	for _, revokedCert := range crl.RevokedCertificateEntries {
		if req.SerialNumber.Cmp(revokedCert.SerialNumber) == 0 {
			ocspResponseTemplate.Status = ocsp.Revoked
			ocspResponseTemplate.RevokedAt = revokedCert.RevocationTime
			ocspResponseTemplate.RevocationReason = revokedCert.ReasonCode
			break
		}
	}

	ocspResponseDER, err := ocsp.CreateResponse(caCert, caCert, ocspResponseTemplate, caPrvKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create OCSP response")
	}

	return ocspResponseDER, nil
}
