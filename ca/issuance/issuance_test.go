package issuance

import (
	"fmt"
	"testing"

	"github.com/agclqq/prow-framework/ca/csr"
	"github.com/agclqq/prow-framework/ca/prvkey"
)

var (
	caKey  = []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA5AVmioMrASXdzgF0grAj73LIsKGSPvvnfH2jBmyR+TgOpsV1\nKJfofqDN6DAQ+qbKFCDfe62Kfcw5ph7q02eQ5YzfsNGK0ez+CmzmPQgt0gS2sW7E\ndfK8joeJp9XMZiVtCQ7UbCgsKapl1+nFaBKBnW4AWXF2gN+cn1FdulsPpygR/px8\npPgKiaHDbS6OgMUbJmwAcOGs5v/QJUOSLDmcawBwP++PBxZewTnrQvyj7xXmbFzI\nVLnlVZEfwUSEL9o+WubtmdsAI6jyztD0/NhdTEKQ2U8GMkEXgJwNXcLPPKjh6WVg\nHkQBXlYXXjEwVwbogTJNqMr7yLS8mr3e9qMgcwIDAQABAoIBADeiMbCd7EItcP6w\nUDMKstnbUaf24+3GHGa9aKdKmhsKWFjMWJd86NbseRCrmZjuVOOwWZadcurahz5G\n0NifrjSzuGg11/78Kcd1Zn+BnVxelgyYkAqPHP5Rh36RpXtOqlnJan6xFoVb89lI\nSkfoLAOzMRahnl43MMmWWp37VchcgnhIXCH/qug6OFN2kDtwBb1fU9o5EA9NX69V\nPabae1AaYwd2KsVYfu5kYKxB56a5pj9gwUkhLY8XxAHealahHWyxmS3PljSSbcDX\n2x6E39AwMrgDB3BXn76KS2NJm/S8LPJTbu7cc8qwjeZj2VxqludEBAZ36C2JEult\nJNgkMJECgYEA5+mc8aYtqXJmvfpdM1H2ePrk+57O8T6fRGFe2uwvsY5lsvp68VZn\ngvZiV6yKTh3W6/+/dD3ufez23v1/CSbQlinxEma7l4EjPXNU3GuJO87qjWUGrZPy\noDboU96EdOWgBmrwWj4+7YhyzUoHnxqpZObSm47Rkx84xeJ/imHTChcCgYEA+7RR\nBviU9JOcOaVtXNDbOwLgfsKPkxUgTEpMceSu+lmdifpJEKZaMXa0XpWa141oI/AS\nm251OiDCrTjtBaIaJX0WbDIeKOFmQnL+YtDmpOot/RxkRsLja/M/VT/DkXkR3DC9\nDD4KWa9UiOMMNX6E+T8yuAynrvR0fvTz2PamQgUCgYB8EnnKtrM7Ml9RSD7QlAsf\nEmurSn1Ah9ZBiS5sRWwGvD4gkO1xbF6YrCRU75RW0pQHUp4lHHUZncs95bUvOjrh\n+7Jju96k4Yvu9mLyQf37p2nJF8GI39wwZu/I3wVSXP9OL6xDO4YDIrr4paCKOINj\n3jHS04fABDYleFXvvQJhJQKBgQDDYKEwrw44lsfCe2VbkYdK1B3cZzu5KMHsVhPm\nqGMmUx+VNaE3elkyYfj6HliWDt6SXsyit+fo2fsjKLfbEowHI4SfMXv1sZiF5esO\nWyddRaWy/jHcN3T+m09C5f7xUbAKYg6sjQ/Ns+oDY3Jbp7yiGGtPMAuNI5W14n/R\nnwtI3QKBgGA0lXFND1noQoortQTKXUlp+dq9U0Emvn6J4ilbNn+FGYcOR5JLXfP9\nfI6yw6l/1eNsT4WHES5xrxK7o6nGt3Wi1jmdQNewtqCAN9u2HAAAsuKil76No6Xy\nt6v2+xt0KKfn1QlhJt3TM28CykM7iZHXhA3RZoKdPC508trE/7OD\n-----END RSA PRIVATE KEY-----")
	caCert = []byte("-----BEGIN CERTIFICATE-----\nMIIDzDCCArSgAwIBAgIRAN+0HUzkiF8KMijV5V3md0kwDQYJKoZIhvcNAQELBQAw\ncDELMAkGA1UEBhMCQ04xEDAOBgNVBAgTB0JlaWppbmcxEDAOBgNVBAcTB0JlaWpp\nbmcxEzARBgNVBAoMCm15X2NvbXBhbnkxFjAUBgNVBAsMDW15X2RlcGFydG1lbnQx\nEDAOBgNVBAMMB215X25hbWUwHhcNMjQwNjE3MDc1NzU4WhcNMjUwNjE3MDc1NzU4\nWjBwMQswCQYDVQQGEwJDTjEQMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVp\namluZzETMBEGA1UECgwKbXlfY29tcGFueTEWMBQGA1UECwwNbXlfZGVwYXJ0bWVu\ndDEQMA4GA1UEAwwHbXlfbmFtZTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC\nggEBAOQFZoqDKwEl3c4BdIKwI+9yyLChkj7753x9owZskfk4DqbFdSiX6H6gzegw\nEPqmyhQg33utin3MOaYe6tNnkOWM37DRitHs/gps5j0ILdIEtrFuxHXyvI6HiafV\nzGYlbQkO1GwoLCmqZdfpxWgSgZ1uAFlxdoDfnJ9RXbpbD6coEf6cfKT4Comhw20u\njoDFGyZsAHDhrOb/0CVDkiw5nGsAcD/vjwcWXsE560L8o+8V5mxcyFS55VWRH8FE\nhC/aPlrm7ZnbACOo8s7Q9PzYXUxCkNlPBjJBF4CcDV3Czzyo4ellYB5EAV5WF14x\nMFcG6IEyTajK+8i0vJq93vajIHMCAwEAAaNhMF8wDgYDVR0PAQH/BAQDAgKEMB0G\nA1UdJQQWMBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MB0G\nA1UdDgQWBBTRL8qQ/Qm7PforiDXCEW5YrLxIfTANBgkqhkiG9w0BAQsFAAOCAQEA\nSXCZeKRo7tIQo24ddtGwF8W9sIPorRAKxxg6IyVdYikq1gKEOg11TvIXhB1Ivo1C\n5xb/Q6DI3+jIR/bJEqLuUcEm/0wclQuZYykR74Mr7RoCjs2d4LG9RwBXsS9XyFff\nEAupynVbEm96kN7MIsIbWk2aG83Tg9d/lN0X3Xspq1jS1WswPq14KoHmDn6RecqW\nlg+E9jUGyOQhmqfa4GMqY+/JXLDdeEXtqbgU+ei49tOEE8nXiWn4XFgi3sOS9PFc\nuDxTBBWK+g9Ka5QmgccENz7ZqpBwtms33baf1Mgldz9HY7xaG4qsrUFVL89dgpcl\nH9+vOoItAhxQ8LJAyBF4rg==\n-----END CERTIFICATE-----")
)

func TestCert_Sign(t *testing.T) {
	type fields struct {
		cert *Cert
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "t1", fields: fields{cert: NewCert(caKey, caCert, func() []byte {
			key, err := prvkey.Gen(2048)
			if err != nil {
				return nil
			}
			csrPem, err := csr.NewCsr(key, []string{"CN"}, []string{"BJ"}, []string{"bj"}, []string{"my_company"}, []string{"my_department"}, "t1", csr.WithDns([]string{"localhost"})).Gen()
			if err != nil {
				return nil
			}
			return csrPem
		}(), WithDays(10), WithIssueType(IssueTypeServer))}, wantErr: false},
		{name: "t2", fields: fields{cert: NewCert(caKey, caCert, func() []byte {
			key, err := prvkey.Gen(2048)
			if err != nil {
				return nil
			}
			csrPem, err := csr.NewCsr(key, []string{"CN"}, []string{"BJ"}, []string{"bj2"}, []string{"my_company2"}, []string{"my_department2"}, "t2").Gen()
			if err != nil {
				return nil
			}
			return csrPem
		}())}, wantErr: false},
		{name: "t3", fields: fields{cert: NewCert(caKey, caCert, func() []byte {
			// openssl genpkey -algorithm RSA  -pkeyopt rsa_keygen_bits:2048
			key := []byte("-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDF16X4mtF7zM26\nHAcsQHhaT8ffHTJRH4/fJDU7GCF9c9o/Q6oNprnsYcC6f4ErkuailMx/jKjomkMP\nPBPDPOtVM2AEykpnWvBQ8j33Ea08KyC8UMAYEfe7WCXtqPfsrc560fforKPENWMF\nSLnRPaDxNv8TuR1QpfzQ7Oxp2WgJ9uuCt1c7z5Hjn3DGBfqtV6AcYUAYHnxo1zvn\n1bsQcSWeNuEEsqZ38cHn57IPvNKhcq/smXm7BqUXnspjBrFOdqd+IF1YDorzQjhS\nQyK/c9QovMMZOqFQJjczIwHqC58ju5Vd453xKqXfvHnpFUAqpx+pLJAodXMqpevo\nrjz0LIYBAgMBAAECggEAG3V6y65xNMWQKCyLslCgY6h/DTB4M1o2FbpoyPLocwkJ\nWY6Co7JoS66lmTzpKKsS563PVESpjN8cP5kPBSIHZ6Phx8hr2zx77kAw6YHCkX9K\n49gxUSXtREtPuFSjVG4rIlDSH7EWab0fKTSW1bvAArqnXI1szCy9kiHQDkDmd7tX\nK+CK1op/RGWud+5Ob+UcZb9gmIOKsf5EnOfMd0TkHRH4sMkM74tLBtod2t26Fu6W\nunvBK8aFVtUH3DZnmqhJmG8ZazD/8UJ4zkluJZCQIz5GV9xB09VwvqqNH2Y/8Gfw\nKUYdda+6YplS6mqSzjLqOuywXVRxX5QiK8k83UrQIQKBgQDh9bmvqtri1TshJxeF\ncDTwTOZbUHLbwoU9v4VU3JcvJ1vbkizaa2BgzeN1UAAfxeUxJ5bG3y4TiucT8Rju\nTPslm4Kfjj68PAuhlLM1QeHsLUcKNeaSKKEAs6zIT9fbiQA6yjzAoBofcKj9Iwp4\n0vTDDk/G++U4ELbvsyxM3VgQJwKBgQDgJPpQO3jpLVLkvWsDljJTpF6SJstIrxL9\nKuRJaLlyBMGkVBPu7mQNDxB9ALStxes5nf0xLu+6YIVQo/Y1zezdgKeYks5eTN9l\naREKkEtFw6qWt5NHqlr68RJgDHxIQgZQaSVvyOzy5ICIJrQsMJWHq43jr3e5v/Ev\n/yZIKm1plwKBgQCK9D+CRcFhaNt54b5XMs97Tu8CDJD1j8O8W0C1FQpr1vpoJpYq\no4mbPkG4bMAGyf3NopjYJ3sATZUY8FTyhqiTUfScBi+SNiK49ObXw3IZeSaMouTt\n0Mph0hxY+rC1sqRPgvqlQk+OMgvZz2irMJ+QLAbnSRSGy9CTy01c32k+VQKBgQDG\nlRxMHxS1idlKHOOFzvkRj6vV0pcB81JgiDKvMyAxezNQcskiQ4TS6QjTpt9sodAQ\nQQAEJjBwMHmMg2dsLeBwMj7J9y7s7zBw+VAGyuZVjdBCLaxHrw9iClkcTZOCtTRA\n45cuXZZIb9fMSHYSPI0OIRjZoyjwobR+sJBrGWPMSwKBgHIL9ypcAK+kW/2IcTfB\n7woa4oH30M/gSeLX9FECXYEa3x+kO/HTfZbjnqq9GjYdT+3DLfmsdBczHh2jaGKS\nrC40hBW3/jKlbuQ+wTMnpgLcFSILP64dMRfOagIO7SE/oalcwhVBZ0bci82EL6/A\nk1xdYiPTpycZA+6DOBDOyyyT\n-----END PRIVATE KEY-----")
			csrPem, err := csr.NewCsr(key, []string{"CN"}, []string{"BJ"}, []string{"bj3"}, []string{"my_company3"}, []string{"my_department3"}, "t3").Gen()
			if err != nil {
				return nil
			}
			return csrPem
		}())}, wantErr: false},
		{name: "t4", fields: fields{cert: NewCert(caKey, caCert, func() []byte {
			key, err := prvkey.Gen(2048)
			if err != nil {
				return nil
			}
			csrPem, err := csr.NewCsr(key, []string{"CN"}, []string{"BJ"}, []string{"bj4"}, []string{"my_company4"}, []string{"my_department4"}, "t4").Gen()
			if err != nil {
				return nil
			}
			return csrPem
		}(), WithDays(10), WithIssueType(IssueTypeIntermediate))}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.cert.Sign()
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(string(got))
		})
	}
}
