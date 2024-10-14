package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/agclqq/prow-framework/ca/csr"
	"github.com/agclqq/prow-framework/ca/info"
	"github.com/agclqq/prow-framework/ca/issuance"
	"github.com/agclqq/prow-framework/ca/ocsp"
	"github.com/agclqq/prow-framework/ca/prvkey"
	"github.com/agclqq/prow-framework/ca/revoke"
	"github.com/agclqq/prow-framework/ca/selfsign"
)

var (
	shutdown     = make(chan struct{}, 1)
	shutdownOcsp = make(chan struct{}, 1)
	caKey        []byte
	caCert       []byte
	interKey     []byte
	interCert    []byte
	svrPrvKey    []byte
	svrCert      []byte
	cliPrvKey    []byte
	cliCert      []byte
	crl          []byte
)

// 创建ca证书
func createCa() ([]byte, []byte, error) {
	return selfsign.NewCa([]string{"CN"}, []string{"Beijing"}, []string{"beijing"}, []string{"my_company"}, []string{"my_department"}, "im_ca").Sign()
}

// 创建server证书
func createServerCert(caCert, caKey []byte, ocsp []string) ([]byte, []byte, error) {
	prvKey, err := prvkey.Gen(2048)
	if err != nil {
		return nil, nil, err
	}
	csrByte, err := csr.NewCsr(prvKey, []string{"CN"}, []string{"Beijing"}, []string{"beijing"}, []string{"my_company"}, []string{"my_department"}, "im_server", csr.WithDns([]string{"localhost"})).Gen()
	if err != nil {
		return nil, nil, err
	}
	cert, err := issuance.NewCert(caCert, caKey, csrByte, issuance.WithIssueType(issuance.IssueTypeServer), issuance.WithOcspServer(ocsp)).Sign()
	if err != nil {
		return nil, nil, err
	}
	return cert, prvKey, nil
}

// 创建client证书
func createClientCert(caCert, caKey []byte, ocsp []string) ([]byte, []byte, error) {
	prvKey, err := prvkey.Gen(2048)
	if err != nil {
		return nil, nil, err
	}
	csrByte, err := csr.NewCsr(prvKey, []string{"CN"}, []string{"Beijing"}, []string{"beijing"}, []string{"my_company"}, []string{"my_department"}, "im_client").Gen()
	if err != nil {
		return nil, nil, err
	}
	cert, err := issuance.NewCert(caCert, caKey, csrByte, issuance.WithOcspServer(ocsp)).Sign()
	if err != nil {
		return nil, nil, err
	}
	return cert, prvKey, nil
}

// 创建中间证书
func createInterCert(caCert, caKey []byte, ocsp []string) ([]byte, []byte, error) {
	prvKey, err := prvkey.Gen(2048)
	if err != nil {
		return nil, nil, err
	}
	csrByte, err := csr.NewCsr(prvKey, []string{"CN"}, []string{"Beijing"}, []string{"beijing"}, []string{"my_company"}, []string{"my_department"}, "im_inter").Gen()
	if err != nil {
		return nil, nil, err
	}
	cert, err := issuance.NewCert(caCert, caKey, csrByte, issuance.WithIssueType(issuance.IssueTypeIntermediate), issuance.WithOcspServer(ocsp)).Sign()
	if err != nil {
		return nil, nil, err
	}
	return cert, prvKey, nil
}

// 吊销证书
func revokeCert(cert []byte) ([]byte, error) {
	return revoke.NewRevoke(caCert, caKey, cert).Revoke()
}

func createKeyCert() error {
	var err error
	caCert, caKey, err = createCa()
	if err != nil {
		return err
	}

	svrCert, svrPrvKey, err = createServerCert(caCert, caKey, nil)
	if err != nil {
		return err
	}

	cliCert, cliPrvKey, err = createClientCert(caCert, caKey, []string{"https://localhost:8081/ocsp"})
	return err
}
func TestTLSConfig(t *testing.T) {
	err := createKeyCert()
	if err != nil {
		t.Error(err)
		return
	}
	//go svr(t) //开启server

	time.Sleep(1 * time.Second)

	type fields struct {
		crlFunc func() error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "t1", fields: fields{crlFunc: func() error {
			//吊销一个无关的证书
			tmpCert, _, err := createClientCert(caCert, caKey, []string{"https://localhost:8081/ocsp"})
			if err != nil {
				return err
			}
			infos, err := info.NewCert(tmpCert).GetInfo()
			if err != nil {
				return err
			}
			t.Logf("tmpCert ,revoke SerialNumber %d", infos.SerialNumber)
			crl, err = revokeCert(tmpCert)
			if err != nil {
				return err
			}
			return nil
		}}, wantErr: false},
		{name: "t2", fields: fields{func() error {
			//把客户端的证书吊销
			infos, err := info.NewCert(cliCert).GetInfo()
			if err != nil {
				return err
			}
			t.Logf("cliCert ,revoke SerialNumber %d", infos.SerialNumber)
			crl, err = revokeCert(cliCert)
			if err != nil {
				return err
			}
			return nil
		}}, wantErr: true},
		{name: "t3", fields: fields{func() error {
			//重新生成客户端证书
			cliCert, cliPrvKey, err = createClientCert(caCert, caKey, []string{"https://localhost:8081/ocsp"})
			if err != nil {
				return err
			}
			//吊销无ocsp服务的服务端的证书，此时客户端无法验证，服务可正常访问
			infos, err := info.NewCert(svrCert).GetInfo()
			if err != nil {
				return err
			}
			t.Logf("svrCert ,revoke SerialNumber %d", infos.SerialNumber)
			crl, err = revokeCert(svrCert)
			if err != nil {
				return err
			}
			return nil
		}}, wantErr: false},
		{name: "t4", fields: fields{func() error {
			//重新生成客户端证书
			cliCert, cliPrvKey, err = createClientCert(caCert, caKey, []string{"https://localhost:8081/ocsp"})
			if err != nil {
				return err
			}
			//吊销有ocsp服务的服务端的证书，此时客户端可以验证，服务不可访问
			svrCert, svrPrvKey, err = createServerCert(caCert, caKey, []string{"https://localhost:8081/ocsp"})
			if err != nil {
				return err
			}
			infos, err := info.NewCert(svrCert).GetInfo()
			if err != nil {
				return err
			}
			t.Logf("svrCert ,revoke SerialNumber %d", infos.SerialNumber)
			crl, err = revokeCert(svrCert)
			if err != nil {
				return err
			}
			return nil

		}}, wantErr: true},
		{name: "t5", fields: fields{func() error {
			//使用ca生成中间证书
			interCert, interKey, err = createInterCert(caCert, caKey, []string{"https://localhost:8081/ocsp"})
			if err != nil {
				return err
			}

			//重新生成客户端证书
			cliCert, cliPrvKey, err = createClientCert(interCert, interKey, []string{"https://localhost:8081/ocsp"})
			if err != nil {
				return err
			}
			//重新生成服务端证书
			svrCert, svrPrvKey, err = createServerCert(interCert, interKey, []string{"https://localhost:8081/ocsp"})
			if err != nil {
				return err
			}
			infos, err := info.NewCert(svrCert).GetInfo()
			if err != nil {
				return err
			}
			t.Logf("svrCert SerialNumber %d", infos.SerialNumber)
			//吊销有ocsp服务的服务端的证书
			crl, err = revokeCert(svrCert)
			if err != nil {
				return err
			}
			return nil

		}}, wantErr: true},
		{name: "t6", fields: fields{func() error {
			//使用ca生成中间证书
			interCert, interKey, err = createInterCert(caCert, caKey, []string{"https://localhost:8081/ocsp"})
			if err != nil {
				return err
			}

			//重新生成客户端证书
			cliCert, cliPrvKey, err = createClientCert(interCert, interKey, []string{"https://localhost:8081/ocsp"})
			if err != nil {
				return err
			}

			infos, err := info.NewCert(cliCert).GetInfo()
			if err != nil {
				return err
			}
			t.Logf("cliCert SerialNumber %d", infos.SerialNumber)

			//吊销客户端证书
			crl, err = revokeCert(cliCert)
			if err != nil {
				return err
			}
			//重新生成服务端证书
			svrCert, svrPrvKey, err = createServerCert(interCert, interKey, []string{"https://localhost:8081/ocsp"})
			if err != nil {
				return err
			}
			return nil

		}}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = tt.fields.crlFunc()
			if err != nil {
				t.Error(err)
				return
			}
			svrPort := svr(t) //开启server

			go ocspSvr(t) //开启ocsp server

			time.Sleep(1 * time.Second)

			pair, err := tls.X509KeyPair(cliCert, cliPrvKey)
			if err != nil {
				return
			}

			certPool := x509.NewCertPool()
			if ok := certPool.AppendCertsFromPEM(caCert); !ok {
				t.Error("failed to parse ca certificate")
				return
			}
			if len(interCert) > 0 {
				if ok := certPool.AppendCertsFromPEM(interCert); !ok {
					t.Error("failed to parse inter certificate")
					return
				}
			}

			cli := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						Certificates:          []tls.Certificate{pair},        // 设置客户端证书和私钥
						RootCAs:               certPool,                       // 设置服务端CA证书池用于验证服务端证书
						ClientAuth:            tls.RequireAndVerifyClientCert, // 启用客户端证书验证（即双向证书验证）
						VerifyPeerCertificate: Vpc,
					},
				},
			}
			_, err = cli.Get("http://localhost:" + svrPort)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				time.Sleep(1 * time.Second)
				shutdown <- struct{}{}
				shutdownOcsp <- struct{}{}
				return
			}
			time.Sleep(1 * time.Second)
			shutdownOcsp <- struct{}{}
			shutdown <- struct{}{}
		})
	}
	return
}
func ocspSvr(t *testing.T) {
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		t.Error("failed to parse CA certificate")
		return
	}
	if len(interCert) > 0 {
		if ok := caCertPool.AppendCertsFromPEM(interCert); !ok {
			t.Error("failed to parse inter certificate")
			return
		}
	}

	pair, err := tls.X509KeyPair(svrCert, svrPrvKey)
	if err != nil {
		t.Error(err)
		return
	}
	server := &http.Server{
		Addr: ":8081",
		Handler: func() *http.ServeMux {
			mux := http.NewServeMux()
			mux.HandleFunc("/ocsp", func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					w.Write([]byte("failed to read body"))
					return
				}
				bytes, err := ocsp.Ocsp(body, caCert, caKey, crl)
				if err != nil {
					w.Write([]byte(err.Error()))
					return
				}
				w.Write(bytes)
			})
			return mux
		}(),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{pair},
			ClientCAs:    caCertPool,
		},
		IdleTimeout: 75 * time.Second,
	}

	go func() {
		err = server.ListenAndServeTLS("", "")
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			t.Error(err)
			return
		}
	}()
	<-shutdownOcsp
	server.Shutdown(nil)
}
func svr(t *testing.T) string {
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		t.Error("failed to parse CA certificate")
		return ""
	}
	if len(interCert) > 0 {
		if ok := caCertPool.AppendCertsFromPEM(interCert); !ok {
			t.Error("failed to parse inter certificate")
			return ""
		}
	}

	mux := http.NewServeMux()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	pair, err := tls.X509KeyPair(svrCert, svrPrvKey)
	if err != nil {
		t.Error(err)
		return ""
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{pair},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,

		VerifyPeerCertificate: Vpc,
	}

	server := httptest.NewUnstartedServer(mux)
	server.TLS = tlsConfig
	server.StartTLS()
	urlParse, err := url.Parse(server.URL)
	if err != nil {
		return ""
	}
	return urlParse.Port()
}
