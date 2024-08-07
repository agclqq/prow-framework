package revoke

import (
	"fmt"
	"testing"
)

var (
	caKey  = []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAw0XxSOV/nGs8ijEWAKUbDOWPdTmK3GjOss3LnZD7BiTskWIg\n+jx9B1OTDmQQb1nmX24ZQEjRWpwzrmgQ7VgbM+XjFXhLl2H7G3mZA1QzLT2Ze9mi\n0JGa9OACWbJy55+40OEXX8pAchw10yTi5Ys/qhKVRM2HgLscI4fs/CWREXHuumrH\nT5f4WEEbHENSPUILzEUnQSaWdFRgcwTdVDRGJK6Xsz14/YOSeNweG1ZA6yUMHifn\nZFYbiO17raSrcV4JTzvLeY5TYpgriyxaNBGVNizsfl9cLpswDwlZxpclF7uGfIWs\nIRsomCOosjxr9/p0FAoK26ph4DfR8EixG6PawwIDAQABAoIBAGaxBJfiYT7AQmEm\nOTzzlwssOkpajYUl9PWhNmBRm0F675HxOgh/AP12XRKnWuFENNugydS9tqNhG+iv\nP5+hwwSC8+4Zih89XtHvG6HdiOBU0b+JD4+B0yzOFU5Ygwb+PzJR/XnZohSgc0nr\nzwsKNNva0/cP4x+2xrCEzgM4OlciV2AncQVhj3U0D6g20Vx92kA750QI6phSXTJl\nBrVqGpcdQdSzb4QG0qDO19T6yRn6YEvK8Zhnxe5GoAgbsD2m3+XLwc0AJnpgpR9C\nDcjpbM8V922emplZHClziatIraFMNcJxlZiDSLyCiW8Fg81F/vHp7nPwgIXDrikh\n+g/E/hkCgYEAxKgTYfehnV+t0mLHkvlB1eHVIkECBA1dXHwWpyYigH5RpIUcbojH\nJXFps+56jTPclNsz7AN03PB9bMntcPI5aBjSyGXUVgIUXvjbLix+IUsXOmYBIvbE\nc32F5ADcjh6Q9VTZXjkz+bLdOiG6N16avp/WubVBCiQTGPFYhiv2v2cCgYEA/jMA\nvdcMUC4ceh/9xFYnIHxg2wqOKUJzNHKhwZq5nHYAbAJJvmRCNB6te8XnakORoL6j\nsRNHVH4I/MSBj1PLHs/k2CN7NrGQ/LoTsQo3/v48P+zXPGt8lM9+V2TWdab6hBeN\nWEglcuQ1/pnJTlDZnr2DljcFPAb0er6ghIhVHEUCgYBTubM23HUESX35unB5lIGB\nC/rv8HlpPD9pZrNGSqgZyK39u2ZVcQpIWLbGElw+zbu17HV4oCgbAJCFxpq+oYHr\nXdYv15rFW3FM1eqLCApTJmMnS9JkDmepO+HLJsq//yd8K7m3seb9AjfJzh44AKEl\nU2vZ+N9N7/npfqdPyFvvNQKBgQDNVc+0ieDlZ8oTEIKBtYIXqMDoT0d0pru+0xY8\n+MoUS/GdTd/Zzsz3owxKHhwH55rcOQKrSEJnSwPhgq6RY4OBWTenLEocbSUMMRc0\n/GctMJrknGFk6gKRhmatG8Rs9zwHtaq0dFrjytqe1gUZoQ+ZPcbscXdl/MxB1nh9\ndk8h7QKBgFP9jGOqJPbnUhHQLvcPFG64PE/hWWqugXbr8Zv0/YC/P6hQAUhIw3U4\nTbm+//08p80EIkfYVaR9yCycU4+1+p2uLjOTfGVTtEYLxxxi/RE8jr0eyX5CVB8G\n5yYt7opReHlm+uz1mMVnih5ZtAG7iycfskKIpKS3TWsw/W+l1U9W\n-----END RSA PRIVATE KEY-----")
	caCert = []byte("-----BEGIN CERTIFICATE-----\nMIIDzDCCArSgAwIBAgIRAJ3CRqZa1JdWAjlaSXJYS58wDQYJKoZIhvcNAQELBQAw\ncDELMAkGA1UEBhMCQ04xEDAOBgNVBAgTB0JlaWppbmcxEDAOBgNVBAcTB0JlaWpp\nbmcxEzARBgNVBAoMCm15X2NvbXBhbnkxFjAUBgNVBAsMDW15X2RlcGFydG1lbnQx\nEDAOBgNVBAMMB215X25hbWUwHhcNMjQwNjI0MTEyMzM2WhcNMjUwNjI0MTEyMzM2\nWjBwMQswCQYDVQQGEwJDTjEQMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVp\namluZzETMBEGA1UECgwKbXlfY29tcGFueTEWMBQGA1UECwwNbXlfZGVwYXJ0bWVu\ndDEQMA4GA1UEAwwHbXlfbmFtZTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC\nggEBAMNF8Ujlf5xrPIoxFgClGwzlj3U5itxozrLNy52Q+wYk7JFiIPo8fQdTkw5k\nEG9Z5l9uGUBI0VqcM65oEO1YGzPl4xV4S5dh+xt5mQNUMy09mXvZotCRmvTgAlmy\ncuefuNDhF1/KQHIcNdMk4uWLP6oSlUTNh4C7HCOH7PwlkRFx7rpqx0+X+FhBGxxD\nUj1CC8xFJ0EmlnRUYHME3VQ0RiSul7M9eP2DknjcHhtWQOslDB4n52RWG4jte62k\nq3FeCU87y3mOU2KYK4ssWjQRlTYs7H5fXC6bMA8JWcaXJRe7hnyFrCEbKJgjqLI8\na/f6dBQKCtuqYeA30fBIsRuj2sMCAwEAAaNhMF8wDgYDVR0PAQH/BAQDAgGGMB0G\nA1UdJQQWMBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MB0G\nA1UdDgQWBBTUABsZokckiRHbgnra7U/sAoDvPjANBgkqhkiG9w0BAQsFAAOCAQEA\noZayrJ8KUjz1vMoLSZPHrmlKTnCE4aflyMsaVkIDUyk7lmWhUZ2l7/JcCEScZgOu\nQhiskNcJpSsg0NY/0JHsNZePYOnocL+Xw1vFb+5HgAuCK5q47LnIJHlpE83/VPYY\nn8qGj6VM5v7PxLAZa5vNZHyi5NWgXKE9foppurdTfMn70r55/J2Yxov1nLMIuZ1A\nziYcIklk/620F2Z6xNp3zWs5FNlh2Cw033FH1/U9aNltERxzyg/oEJfuKJhD8Y0w\n5rUPtfRB8/Hudn1aBIs6buyNFraucqqyyP363lR6647QfI6HAdEzlGYTeJtNFTw5\ntPWfgVaZbY5C0XiCY7tsVw==\n-----END CERTIFICATE-----")
)

func TestRevoke_Revoke(t *testing.T) {
	type fields struct {
		caCert []byte
		caKey  []byte
		cert   []byte
		crl    []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "t1",
			fields: fields{
				caKey:  caKey,
				caCert: caCert,
				cert:   []byte("-----BEGIN CERTIFICATE-----\nMIIDyDCCArCgAwIBAgIRALQF7AUD+JF3X1GvIgTtQMEwDQYJKoZIhvcNAQELBQAw\ncDELMAkGA1UEBhMCQ04xEDAOBgNVBAgTB0JlaWppbmcxEDAOBgNVBAcTB0JlaWpp\nbmcxEzARBgNVBAoMCm15X2NvbXBhbnkxFjAUBgNVBAsMDW15X2RlcGFydG1lbnQx\nEDAOBgNVBAMMB215X25hbWUwHhcNMjQwNjI0MTAyOTAwWhcNMjQwNzA0MTAyOTAw\nWjBhMQswCQYDVQQGEwJDTjELMAkGA1UECBMCQkoxCzAJBgNVBAcTAmJqMRMwEQYD\nVQQKDApteV9jb21wYW55MRYwFAYDVQQLDA1teV9kZXBhcnRtZW50MQswCQYDVQQD\nEwJ0MTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMeb4rEaDV3dv28j\n8Wr6gvrY8Ev3VQ0ZDLqNPU+F52tIckuikp+FZoUQ/CtL0qm0TCUcOBbyeqymco2G\nSs2YT6VDAMyr8IWYgDn5BMlHpSrM0WTtBDzn78EiP8f+xHnBvqZ8HC1DgYycnBnN\nt3jY2S29nVRT+J4b3Pzciyi8rz7TWJSLQAhbS4miIVytQ1uXSyjxVL1gjMPUlenv\n42VTQm5+H1KJkveX+Yh7tsRvOtXiO9PFqHyjtW4u27f6Wt8LkTnki/qzYN46OwJE\nqtWbxX9R7A3fe92gzJrSXy3SF3W+zCatlQ4cqRUkSAwIqVKdizDwTOvFJ7gIFFS+\nMVTkXPMCAwEAAaNsMGowDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUF\nBwMBMAwGA1UdEwEB/wQCMAAwHwYDVR0jBBgwFoAU0S/KkP0Juz36K4g1whFuWKy8\nSH0wFAYDVR0RBA0wC4IJbG9jYWxob3N0MA0GCSqGSIb3DQEBCwUAA4IBAQCut9jj\nzMWkxfVxY834c/BsliCS9prIsO8HDP2PFaBau0pfHAEmKYI3MgbCHHExJkEOnQ4g\nS0M2JmUtoGT6DxDApJs+aTC67u7VngBulWQ6OOer5BMbMnU8n+VzqVSROrZPqRZu\ncUx0ofFirdZJ5HPJazRXlKDJQqHXd3X9U301Chx8pnzeHMEFHvb2ZiPKuFtBXstN\nCmsrmxUhc/xApW1z9YEyoZDEDFwerE5ZJB6xVKgStEz+ces0VMXbaBkDDmvDZipp\nJwRiDVfgIkGIS4GubX0dIOw2a+zlvPSfWEEMzQmEs1X1PQlsbnnYcp7xL6fM64lZ\noh0luyxUcEa+lp9a\n-----END CERTIFICATE-----"),
			},
			wantErr: false,
		},
		{
			name: "t2",
			fields: fields{
				caKey:  caKey,
				caCert: caCert,
				cert:   []byte("-----BEGIN CERTIFICATE-----\nMIIDxzCCAq+gAwIBAgIQSmS99YiZ1+gWGCike/d84TANBgkqhkiG9w0BAQsFADBw\nMQswCQYDVQQGEwJDTjEQMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamlu\nZzETMBEGA1UECgwKbXlfY29tcGFueTEWMBQGA1UECwwNbXlfZGVwYXJ0bWVudDEQ\nMA4GA1UEAwwHbXlfbmFtZTAeFw0yNDA2MjUwMjEwMTlaFw0yNDA3MDUwMjEwMTla\nMGExCzAJBgNVBAYTAkNOMQswCQYDVQQIEwJCSjELMAkGA1UEBxMCYmoxEzARBgNV\nBAoMCm15X2NvbXBhbnkxFjAUBgNVBAsMDW15X2RlcGFydG1lbnQxCzAJBgNVBAMT\nAnQxMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAomhfB/eZext80899\nXZsW5KCMccSzSpp1Y9cexfyb6RHW0Me4ux1LnFju4IKZxyQTmkr6CJDo1exdHXTc\nEjRYvXF+Zv4LLETo0+OFtYfVuVJpcrXe0VFpmDAGGbB62Up9aDiaHLf7fyK1cvmw\nswEbGOt45cUTFsHbn7flcLdG8sInR1BUuxnge9mITALMutnFXqt1wNSljw1EvYcc\nM3eHwC/bnJw0TRHcQ0OQuc/kPneDgM8UwebjV4ayH3fJdoQOHHEWmYHsO8ZuqRx9\nbMVkcs2mHifgrvQTBPH9fEznxZjLRSMlIOH6F2Ih3aP6Nd74NrasvpICxzTM3h7l\nNMbn2QIDAQABo2wwajAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUH\nAwEwDAYDVR0TAQH/BAIwADAfBgNVHSMEGDAWgBTRL8qQ/Qm7PforiDXCEW5YrLxI\nfTAUBgNVHREEDTALgglsb2NhbGhvc3QwDQYJKoZIhvcNAQELBQADggEBAH7z1upd\nAdmnv1UF5+KMPHaT/gU3gYs9fLGTFJv2l23hEjGp4puQidjYcEDpheDhbGOedVkj\nFaklL5bNaaRX1KilKCxD2tMTTtoaXtl2AJEgNCd3vbE+IBkqV3BIDA4Ki7HGKCxy\n6dl2CrD6G9He90UHhF6VNVNDuw8c2udA7rFOa7vp3cm65EQuElU6F/9kvWab46jE\nwze8d5qcvjzkafWamKGhCH8c9rhLDhk6UM0T0XtODmeG2H7KJr10Qy2a7GCHGP2y\n9QGOcQ1dYKM3iqQ16RsfcJLLScOGKwph6ndgJpCE4N1FgJ57+vB2Cdvp5uTkGXMp\nGl+OCMH++qQAWMI=\n-----END CERTIFICATE-----"),
				crl:    []byte("-----BEGIN X509 CRL-----\nMIICEDCB+QIBATANBgkqhkiG9w0BAQsFADBwMQswCQYDVQQGEwJDTjEQMA4GA1UE\nCBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzETMBEGA1UECgwKbXlfY29tcGFu\neTEWMBQGA1UECwwNbXlfZGVwYXJ0bWVudDEQMA4GA1UEAwwHbXlfbmFtZRcNMjQw\nNjI1MDE1NjUyWhcNMjQwNzAyMDE1NjUzWjAkMCICEQC0BewFA/iRd19RryIE7UDB\nFw0yNDA2MjUwMTU2NDdaoC8wLTAfBgNVHSMEGDAWgBTUABsZokckiRHbgnra7U/s\nAoDvPjAKBgNVHRQEAwIBATANBgkqhkiG9w0BAQsFAAOCAQEAWZZap4jbxhNae1LL\nZhgdbi0t5UbRl2f+W5albOO+1dm+y83IgycWRNLuTqoZd088tH4s6gP+QUwsIA7u\ngErTTQzidMv99G4iekbg9ndVH9cPv44/Orsnv2V6mr1tnRa6S8jqIieMK2vrxGpQ\nY2D/MzrNp+auq/oLW66pJ6eE1Inp2N/LWORzwj5qjiQIx6OzlnYTONb/6tWmBkDX\n3Kbxfc6fi1amCgsBnJzcdRuek6ePW4EC920OVUh0xzfE3e5w/wZhKzJl9Usd7+Np\ncznwF2VQQpEwCSyVTC02AOTpZ9kun4pqxfZjXEXVR0nLC0cF720fHojed0PHdAKo\n6PuxQw==\n-----END X509 CRL-----"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Revoke{
				caCert: tt.fields.caCert,
				caKey:  tt.fields.caKey,
				cert:   tt.fields.cert,
				crl:    tt.fields.crl,
			}
			got, err := r.Revoke()
			if (err != nil) != tt.wantErr {
				t.Errorf("Revoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(string(got))
		})
	}
}
