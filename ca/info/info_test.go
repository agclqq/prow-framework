package info

import (
	"testing"
)

func TestCert_Get(t *testing.T) {
	type fields struct {
		cert *Cert
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "t1", fields: fields{cert: NewCert([]byte("-----BEGIN CERTIFICATE-----\nMIIDyDCCArCgAwIBAgIRAPQBxqKnDt+mXwX9QVO/Z2QwDQYJKoZIhvcNAQELBQAw\ncDELMAkGA1UEBhMCQ04xEDAOBgNVBAgTB0JlaWppbmcxEDAOBgNVBAcTB0JlaWpp\nbmcxEzARBgNVBAoMCm15X2NvbXBhbnkxFjAUBgNVBAsMDW15X2RlcGFydG1lbnQx\nEDAOBgNVBAMMB215X25hbWUwHhcNMjQwNjE5MTAwNjAzWhcNMjQwNjI5MTAwNjAz\nWjBhMQswCQYDVQQGEwJDTjELMAkGA1UECBMCQkoxCzAJBgNVBAcTAmJqMRMwEQYD\nVQQKDApteV9jb21wYW55MRYwFAYDVQQLDA1teV9kZXBhcnRtZW50MQswCQYDVQQD\nEwJ0MTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAM1Gu8HNz5+wLco0\n70i4wjDtfFUsUfg2aZkQTLLlFih1Hnqkc49+fMxcWRdlgNgad0pp/SE24utIcW+M\nXJngq1JBnqrHImp4QsnCPHj7j1TSlRpnwmnN6qFZ/f3nBUJO1fzAM6Lq25It+O/f\npVfPHJ85pBCN8+NOwmwlkd24n4/Xt0nAk01QnAOBR8XGtkFmWH/kX6m75gYMOrqK\n7a01krqzt9H6F7Fi/xvMqCUeBTCFvM1vpazXBbVny8axCHG0ffTcJ2/l93TmeWvh\nmNIukoY7zmrmcTLJe+PW9QuTYac2h9kbyR7jOJyQrvW8gpoN5moQV5aT285eZPzR\nVKqDu6sCAwEAAaNsMGowDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUF\nBwMBMAwGA1UdEwEB/wQCMAAwHwYDVR0jBBgwFoAU0S/KkP0Juz36K4g1whFuWKy8\nSH0wFAYDVR0RBA0wC4IJbG9jYWxob3N0MA0GCSqGSIb3DQEBCwUAA4IBAQCD7/J+\nqSxgkS7eaoOBOHPNiGigXFQ5bdVbIT/axRjd3nBHovA0/urTturRaChJS35Kpq/y\ncYevw+iOFP216fC42HYKfbswILJH5ytUB5MSSTwumszcP+SgoX2hvoMWlmZRVc23\nTTehllKWjOO1RR7y7y/La/Cv3RFNFOeVfGx+pZgOygFPgPsx3I95ioW4KKHdIGrz\nhs5BvXvg5l106oKq47/+1H6a17qcr80SnOphBKAqQ3UdNwyceQ3RPrTB6OVM0scM\nxptOiSbK5o0utOST5v50ms4FQ2+O6pP/zG6HcO9IBamgpB/uWTKF9VDH2/fhgEQ2\ntXl70QyNIBEi2851\n-----END CERTIFICATE-----"))}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.cert.GetInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("cert signature: %v", got.Signature)
			t.Logf("cert Issuer: %v", got.Issuer)
			t.Logf("cert subject: %v", got.Subject)
			t.Logf("cert SerialNumber: %v", got.SerialNumber)
			t.Logf("cert NotBefore: %v", got.NotBefore)
			t.Logf("cert NotAfter: %v", got.NotAfter)
			t.Logf("cert KeyUsage: %v", got.KeyUsage)
			t.Logf("cert ExtKeyUsage: %v", got.ExtKeyUsage)
			t.Logf("cert BasicConstraintsValid: %v", got.BasicConstraintsValid)
			t.Logf("cert IsCA: %v", got.IsCA)
			t.Logf("cert MaxPathLen: %v", got.MaxPathLen)
			t.Logf("cert MaxPathLenZero: %v", got.MaxPathLenZero)
			t.Logf("cert SubjectKeyId: %v", got.SubjectKeyId)
			t.Logf("cert AuthorityKeyId: %v", got.AuthorityKeyId)
			t.Logf("cert DNSNames: %v", got.DNSNames)
		})
	}
}

func TestCrl_Get(t *testing.T) {
	type fields struct {
		crl []byte
	}
	tests := []struct {
		name          string
		fields        fields
		wantRevokeNum int
		wantErr       bool
	}{
		{name: "t1", fields: fields{crl: []byte("-----BEGIN X509 CRL-----\nMIICNDCCARwCAQEwDQYJKoZIhvcNAQELBQAwcDELMAkGA1UEBhMCQ04xEDAOBgNV\nBAgTB0JlaWppbmcxEDAOBgNVBAcTB0JlaWppbmcxEzARBgNVBAoMCm15X2NvbXBh\nbnkxFjAUBgNVBAsMDW15X2RlcGFydG1lbnQxEDAOBgNVBAMMB215X25hbWUXDTI0\nMDYyNTAyMTIwMloXDTI0MDcwMjAyMTIwMlowRzAiAhEAtAXsBQP4kXdfUa8iBO1A\nwRcNMjQwNjI1MDE1NjQ3WjAhAhBKZL31iJnX6BYYKKR793zhFw0yNDA2MjUwMjEx\nMjFaoC8wLTAfBgNVHSMEGDAWgBTUABsZokckiRHbgnra7U/sAoDvPjAKBgNVHRQE\nAwIBAjANBgkqhkiG9w0BAQsFAAOCAQEAPGrfCrE0iDGWYsATZQaOCXbWaRWTaHN1\nCiciE72vzHJ+c7xyLvQkpa3yoHRgm441Rf8VLuZ/BpW7GQO/ivcRxh3KBaFhalhx\ng3R5y0PKwVaHgRQ/IHYtR9xlPhgOjbMWSmZaS3KDgmuut1EMCMypvx5XMKeKyW4f\n3RUManFkU26moRNSQy9e809Y4PpnTJjCYyOEXmDweJ2jGLtXc+M4QWxnr/FWczYq\ncq/030a2KuWP0H+Dy3BzKIfAaJXBFCNr97dwxcTF8rAqrYFS6/vQpQbYOHxrNxdw\n77wCy7PY0l1O2NIUd4slm5gcnkyIORL1dCn8o3mAs2RGDxZMTeQcZA==\n-----END X509 CRL-----")}, wantRevokeNum: 2, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Crl{
				crl: tt.fields.crl,
			}
			got, err := i.GetInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.RevokedCertificateEntries) != tt.wantRevokeNum {
				t.Errorf("GetInfo() got = %v, want %v", len(got.RevokedCertificateEntries), tt.wantRevokeNum)
			}
		})
	}
}
