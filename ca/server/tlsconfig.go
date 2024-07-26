package server

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/crypto/ocsp"
)

func Vpc(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	if len(rawCerts) == 0 {
		return errors.New("no client certificate provided")
	}
	for i, _ := range rawCerts {
		cliCert := verifiedChains[i][0] // 客户端证书
		err := verifyOcsp(cliCert)
		if err != nil {
			return err
		}
	}
	return nil
	//// 解析客户端证书
	//cert, err := x509.ParseCertificate(rawCerts[0])
	//if err != nil {
	//	return fmt.Errorf("failed to parse client certificate: %v", err)
	//}
	//
	//// 验证客户端证书
	//opts := x509.VerifyOptions{
	//	Intermediates: x509.NewCertPool(),
	//}
	//for _, cert := range verifiedChains[0][1:] {
	//	opts.Intermediates.AddCert(cert)
	//}
	//
	//if err = verifyClientCert(cert, opts); err != nil {
	//	return err
	//}
	//
	//return nil
}

func verifyOcsp(cert *x509.Certificate) error {
	ocspReq, err := ocsp.CreateRequest(cert, cert, nil)
	// 检查 OCSP
	if err != nil {
		return fmt.Errorf("failed to create OCSP request: %v", err)
	}
	if len(cert.OCSPServer) == 0 {
		return nil // No OCSP server provided
	}
	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	ocspResp, err := cli.Post(cert.OCSPServer[0], "application/ocsp-request", bytes.NewReader(ocspReq))
	if err != nil {
		return fmt.Errorf("failed to get OCSP response: %v", err)
	}
	defer ocspResp.Body.Close()

	if ocspResp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get OCSP response: %v", ocspResp.Status)
	}

	ocspRespData, err := io.ReadAll(ocspResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read OCSP response: %v", err)
	}

	ocspResult, err := ocsp.ParseResponse(ocspRespData, nil)
	if err != nil {
		return fmt.Errorf("failed to parse OCSP response: %v", err)
	}

	if ocspResult.Status == ocsp.Revoked {
		return fmt.Errorf("certificate with serial number %s has been revoked", cert.SerialNumber)
	}

	return nil
}

// 验证客户端证书
func verifyClientCert(cert *x509.Certificate, opts x509.VerifyOptions) error {
	// 验证证书链
	chains, err := cert.Verify(opts)
	if err != nil {
		return fmt.Errorf("failed to verify certificate chain: %v", err)
	}

	// 检查证书是否被吊销
	for _, chain := range chains {
		for _, cert := range chain {
			if err := checkCertificateRevocation(cert, opts); err != nil {
				return fmt.Errorf("certificate revoked: %v", err)
			}
		}
	}

	return nil
}

// 使用 CRL 或 OCSP 检查证书吊销状态
func checkCertificateRevocation(cert *x509.Certificate, opts x509.VerifyOptions) error {
	// 使用 CRL 检查
	for _, crlURL := range cert.CRLDistributionPoints {
		if err := checkCRL(cert, crlURL); err != nil {
			return err
		}
	}

	// 使用 OCSP 检查
	if err := checkOCSP(cert, opts); err != nil {
		return err
	}

	return nil
}

// 使用 CRL 检查证书是否被吊销
func checkCRL(cert *x509.Certificate, crlURL string) error {
	// 下载并解析 CRL
	resp, err := http.Get(crlURL)
	if err != nil {
		return fmt.Errorf("failed to download CRL: %v", err)
	}
	defer resp.Body.Close()

	crlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read CRL: %v", err)
	}

	crl, err := x509.ParseRevocationList(crlData)
	if err != nil {
		return fmt.Errorf("failed to parse CRL: %v", err)
	}

	// 检查证书是否在 CRL 中
	for _, revoked := range crl.RevokedCertificateEntries {
		if cert.SerialNumber.Cmp(revoked.SerialNumber) == 0 {
			return fmt.Errorf("certificate with serial number %s has been revoked", cert.SerialNumber)
		}
	}

	return nil
}

// 使用 OCSP 检查证书是否被吊销
func checkOCSP(cert *x509.Certificate, opts x509.VerifyOptions) error {
	if len(cert.OCSPServer) == 0 {
		return nil // No OCSP server provided
	}
	//issuer := opts.Intermediates[0] // 假设第一个中间证书是颁发者

	ocspReq, err := ocsp.CreateRequest(cert, nil, nil)
	// 检查 OCSP
	if err != nil {
		return fmt.Errorf("failed to create OCSP request: %v", err)
	}

	ocspResp, err := http.Post(cert.OCSPServer[0], "application/ocsp-request", bytes.NewReader(ocspReq))
	if err != nil {
		return fmt.Errorf("failed to get OCSP response: %v", err)
	}
	defer ocspResp.Body.Close()

	ocspRespData, err := io.ReadAll(ocspResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read OCSP response: %v", err)
	}

	ocspResult, err := ocsp.ParseResponse(ocspRespData, nil)
	if err != nil {
		return fmt.Errorf("failed to parse OCSP response: %v", err)
	}

	if ocspResult.Status == ocsp.Revoked {
		return fmt.Errorf("certificate with serial number %s has been revoked", cert.SerialNumber)
	}

	return nil
}

// 获取颁发者证书
func getIssuerCert(cert *x509.Certificate, opts x509.VerifyOptions) (*x509.Certificate, error) {
	for _, intermediateCert := range getAllCertsFromPool(opts.Intermediates) {
		if cert.CheckSignatureFrom(intermediateCert) == nil {
			return intermediateCert, nil
		}
	}

	for _, rootCert := range getAllCertsFromPool(opts.Roots) {
		if cert.CheckSignatureFrom(rootCert) == nil {
			return rootCert, nil
		}
	}

	return nil, fmt.Errorf("issuer certificate not found")
}

// 从证书池中获取所有证书
func getAllCertsFromPool(pool *x509.CertPool) []*x509.Certificate {
	certs := []*x509.Certificate{}
	for _, certDER := range pool.Subjects() {
		cert, err := x509.ParseCertificate(certDER)
		if err == nil {
			certs = append(certs, cert)
		}
	}
	return certs
}
