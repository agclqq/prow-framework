package main

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"time"

	ocsp1 "golang.org/x/crypto/ocsp"

	"github.com/agclqq/prow-framework/ca/example/ca"
	"github.com/agclqq/prow-framework/ca/example/client"
	"github.com/agclqq/prow-framework/ca/example/server"
	"github.com/agclqq/prow-framework/ca/info"
)

func main() {
	os.Remove("ca.key")
	os.Remove("ca.crt")
	os.Remove("ctl.pem")
	os.Remove("caSvr.key")
	os.Remove("caSvr.crt")
	os.Remove("svr.key")
	os.Remove("svr.crt")
	os.Remove("cli.key")
	os.Remove("cli.crt")
	defer func() {
		os.Remove("ca.key")
		os.Remove("ca.crt")
		os.Remove("ctl.pem")
		os.Remove("caSvr.key")
		os.Remove("caSvr.crt")
		os.Remove("svr.key")
		os.Remove("svr.crt")
		os.Remove("cli.key")
		os.Remove("cli.crt")
	}()
	//启动CA服务
	fmt.Println("init ca cert and caSvr")
	go func() {
		err := ca.Svr("caSvr.crt", "caSvr.key")
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	time.Sleep(2 * time.Second)

	//获取CA证书
	fmt.Println("get ca cert")
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		fmt.Println(err)
		return
	}
	x509CaCert, err := info.NewCert(caCert).GetInfo()
	if err != nil {
		return
	}

	//启动服务端
	go func() {
		fmt.Println("init ca cert and caSvr")
		err = server.Svr(caCert, "svr.crt", "svr.key")
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	time.Sleep(2 * time.Second)

	//启动客户端
	fmt.Println("request svr")
	err = client.Cli(caCert, "cli.crt", "cli.key")
	if err != nil {
		fmt.Println(err)
		return
	}

	//吊销服务端证书
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		fmt.Println("failed to parse CA certificate")
		return
	}
	svrHandler, err := os.Open("svr.crt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer svrHandler.Close()
	svrCert, err := io.ReadAll(svrHandler)
	if err != nil {
		return
	}
	x509SvrCert, err := info.NewCert(svrCert).GetInfo()
	if err != nil {
		return
	}
	block, _ := pem.Decode(svrCert)
	if block == nil || block.Type != "CERTIFICATE" {
		fmt.Println("failed to decode PEM block containing certificate")
		return
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	cli, err := client.GetCli(caCert, nil, nil)

	reRevoke, err := cli.Post("https://127.0.0.1:8080/revoke", "", bytes.NewReader(svrCert))
	if err != nil {
		fmt.Println(err)
		return
	}
	if reRevoke.StatusCode != 200 {
		fmt.Println("failed to revoke certificate")
		return
	}
	if len(cert.OCSPServer) > 0 {
		ocspSvr := cert.OCSPServer[0]
		rsOcsp, err := cli.Post(ocspSvr, "", bytes.NewReader(svrCert))
		if err != nil {
			fmt.Println(err)
			return
		}
		if rsOcsp.StatusCode != 200 {
			fmt.Println("failed to get ocsp")
			return
		}
		all, err := io.ReadAll(rsOcsp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		respOcsp, err := ocsp1.ParseResponse(all, x509CaCert)
		if err != nil {
			return
		}
		switch respOcsp.Status {
		case ocsp1.Good:
			fmt.Println(x509SvrCert.SerialNumber.String(), x509SvrCert.Subject.String(), "good")
		case ocsp1.Revoked:
			fmt.Println(x509SvrCert.SerialNumber.String(), x509SvrCert.Subject.String(), "revoked")
		case ocsp1.Unknown:
			fmt.Println(x509SvrCert.SerialNumber.String(), x509SvrCert.Subject.String(), "unknown")
		default:
			fmt.Println(x509SvrCert.SerialNumber.String(), x509SvrCert.Subject.String(), "unknown err")
		}
	}

	if len(cert.CRLDistributionPoints) > 0 {
		crlSvr := cert.CRLDistributionPoints[0]
		rsCrl, err := cli.Get(crlSvr)
		if err != nil {
			return
		}
		if rsCrl.StatusCode != 200 {
			fmt.Println("failed to get crl")
			return
		}
		crlBytes, err := io.ReadAll(rsCrl.Body)
		if err != nil {
			return
		}
		block, _ := pem.Decode(crlBytes)
		if block == nil || block.Type != "X509 CRL" {
			fmt.Println("failed to decode PEM block containing CRL")
			return
		}
		crl, err := x509.ParseRevocationList(block.Bytes)
		if err != nil {
			return
		}
		for _, revoked := range crl.RevokedCertificateEntries {
			fmt.Println("from crl,revoked certificate: ", revoked.SerialNumber.String())
		}
	}
}
