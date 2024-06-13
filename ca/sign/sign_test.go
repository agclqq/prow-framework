package sign

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/agclqq/prow-framework/execcmd"

	"rms/infra/ca"
	"rms/infra/ca/selfsign"
)

var (
	caKey   = ""
	caPem   = ""
	homeDir = ""
)

func homePath() (string, error) {
	// 获取当前用户
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	// 获取家目录
	return currentUser.HomeDir, nil
}
func selfCa() error {
	home, err := homePath()
	if err != nil {
		return err
	}
	homeDir = home

	keyPath, pemPath, err := selfsign.NewCa("CN", "Beijing", "Beijing", "my company", "my department", "my name", selfsign.WithDir(filepath.Join(homeDir, "/ca/")), selfsign.WithDays(1)).Sign()
	if err != nil {
		return err
	}
	caKey = keyPath
	caPem = pemPath
	return nil
}

func TestNewCert(t *testing.T) {
	if err := selfCa(); err != nil {
		t.Errorf("ca() error = %v", err)
		return
	}
	type args struct {
		caKey string
		caPem string
		opts  []CertOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "t1", args: args{caKey: caKey, caPem: caPem, opts: []CertOption{
			WithCountry("CN"), WithState("Beijing"), WithLocality("Beijing"), WithOrganization("my company"), WithOrganizationalUnit("my department"), WithCommonName("my name"),
			WithDir(filepath.Join(homeDir, "/ca/cert/")), WithDays(100), WithDns("my.dns"), WithServerCaType()}}, wantErr: false},
		{name: "t2", args: args{caKey: caKey, caPem: caPem, opts: []CertOption{
			WithCountry("CN"), WithState("Beijing"), WithLocality("Beijing"), WithOrganization("my company"), WithOrganizationalUnit("my department"), WithCommonName("my name"),
			WithDir(filepath.Join(homeDir, "/ca/cert/")), WithDays(100)}}, wantErr: false},
		{name: "t3", args: args{caKey: caKey, caPem: caPem, opts: []CertOption{
			WithCountry("CN"), WithState("Beijing"), WithLocality("Beijing"), WithOrganization("my company"), WithOrganizationalUnit("my department"), WithCommonName("my name"),
			WithDays(100), WithDns("my.dns"), WithServerCaType()}}, wantErr: false},
		{name: "t4", args: args{caKey: caKey, caPem: caPem, opts: []CertOption{
			WithCountry("CN"), WithState("Beijing"), WithLocality("Beijing"), WithOrganization("my company"), WithOrganizationalUnit("my department"), WithCommonName("my name"),
			WithDays(100), WithPrvKey("-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDCf3LQNUS148v1\ni+yLYb49P4NvJdF1GXecA+N5IH/LFCyyjddC0rG519yP6rHF2BI/4NxRb+r7QmgG\nNrxtvfeN7oPIkfhtKFAAPfMfxj7/g8tbmIv5sXvp26DAioVtrbrU1EXi1jRCe4cZ\n4rHEzFe+mkZ2/pc3CD+4kjG+EqrVlGU1d9cQRM66u2ZYUUPDSn5a300dPes7sjPE\nTVYdwpG1YnxRksCUuyS1IiEHtBXosoF7gK0MGfX7YtY8bHcjYOFCU44sn1bn+FfZ\ng5/SZeOWTY1Jy0dcApVVdSR61yB5EJOkBI7RHMlOjIRKJWBiRAKtS3LPBe4/lnmJ\nTWfJVpsZAgMBAAECggEANm4dTuhBYNetmft9CKqjZxeRrDa8rdUhNH+gFqNCMC5m\nrddlAPXet+ARgRMQigoEXW0LqxyzeXpliyuhQuLxVv6DUcuL5txruw2bLu63baFP\n9UO1FH0XbORCUe/SFFYUnYAESM1iVaKlNdjLoAQBoD0jcCSiY8vCrV/4XLVzqo49\neDbJXHhm/WPyuiY5MGxA6ydGi4TPviQS/ZVtRIxnyYr5nKH4px6aobAFCZHm++dl\nVEgkdmZEK+jw7Pim4I08V7X+0g3IofqIrCtUnB9VXuVFEUl+R6hdENkgFSp6nm9I\nveZsLuVjWjcNuSxGSE6dyCo9Lr2sCFEBlr7CazMFwwKBgQD9VJIi/7EBcx7JYsZH\naGhseH/G3BwHgM5N9k/fyupH51z4dH+tKgl9t2PZDQ8W5a41OjfMmPvCYxDEe4RT\ntcZegAU4DHYA4D+rJC5EbMdX3j7Rpc3gY6GpIXmXTUnJ10qU0bHE+u/ugS5734gM\nG2qo5TM4rnVI5SxknMRlLqYiQwKBgQDEjCkbO6meexjuFS4IlJDpiE/v9mxYiGKW\nVqmqlVihC0tWtt6sQCce783OdwNDD4SctudXLzG8KpY3KaXcLsh/Wpzp302nEx2s\nhiectwRppNQEtfspN5YSgiG99SyVU4g98BnS7Qmn2nmxHmyC1xD8nSNXA5Glp0gS\n7nFOkLP9cwKBgHPi6CcSiMp8+yxs/v9Th9F3Hhy+PCRCjB2l+8wIazwRXrpZsL5q\naIUWC5sTGkADObons7bolOLLprP7PQF+Ogyoy7pkGOc1rmp/1pp+mIJdrKcDDjcD\n3MQeCB1qwcKPthJ2Crhtgqy8c6M/EmFXeWdh0hiv1f9Otwwfmsgemuk5AoGBALTr\nZ7M/uiS9nvcY2+TeDH5LEXoLVTQxZr6IS2lQS+MB6HmLn3DjJJ+fkcxpVMFX+XPg\nERb5xEg200s3tQr2rWw9Vo8ZE/uk5v22B6SD+zXbmaY0dVs9ZZDn5HNcyYsy9wg8\niSjVNLwjqTzWin/txB8j7jHcgScA0qFKh1YQcP3tAoGAWz7rIIFIC4Qe2ggeQqNE\nM1N7hSrCRWYl7FKvAfo2Yg/ybXVk0Jh0MPXTl8qt82VTHnjOZZqmgAv1cJRPl3sR\nA21aP/1X+o7iao0vl0dlkMuedutxfJKtdZCdvZoE9lOlTumHpbkKrXCs00ZlNmMX\nSskXgfMWXgPWXM/IOLHYomQ=\n-----END PRIVATE KEY-----")}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, pem, err := NewCert(tt.args.caKey, tt.args.caPem, tt.args.opts...).Sign()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("key: %s , pem: %s", key, pem)
		})
	}
}

func TestCert_signFromCsr(t *testing.T) {
	if err := selfCa(); err != nil {
		t.Errorf("ca() error = %v", err)
		return
	}
	type fields struct {
		caKey string
		caPem string
		Dir   string
		csr   []byte
		opts  []CertOption
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		want1   string
		wantErr bool
	}{
		{name: "t1", fields: fields{
			caKey: caKey, caPem: caPem, opts: []CertOption{WithDir(filepath.Join(homeDir, "/ca/cert/t1/")), WithCsr(func() []byte {
				subPath := filepath.Join(homeDir, "/ca/cert/t1/")
				os.MkdirAll(subPath, os.ModePerm)
				keyPath := filepath.Join(subPath, "prv.key")
				csrPath := filepath.Join(subPath, "server.csr")
				confFilePath := filepath.Join(subPath, "server.cnf")
				os.WriteFile(confFilePath, []byte(fmt.Sprintf(ca.SignServerCertConf(), "CN", "Beijing", "Beijing", "mycompany", "mydept", "myserver", "test.com")), os.ModePerm)
				execcmd.Command("openssl", "genrsa", "-out", keyPath, "2048")
				execcmd.Command("openssl", "req", "-new", "-sha256", "-out", csrPath, "-key", keyPath, "-config", confFilePath)
				csr, _ := os.ReadFile(csrPath)
				return csr
			}()), WithDays(10), WithDns("my.dns"), WithServerCaType()}}, want: "", want1: filepath.Join(homeDir, "/ca/cert/t1/", "server.pem"), wantErr: false},
		{name: "t2", fields: fields{
			caKey: caKey,
			caPem: caPem,
			opts: []CertOption{
				WithDir(filepath.Join(homeDir, "/ca/cert/t2/")),
				WithCsr(func() []byte {
					subPath := filepath.Join(homeDir, "/ca/cert/t2/")
					os.MkdirAll(subPath, os.ModePerm)
					keyPath := filepath.Join(subPath, "prv.key")
					csrPath := filepath.Join(subPath, "client.csr")
					confFilePath := filepath.Join(subPath, "client.cnf")
					os.WriteFile(confFilePath, []byte(fmt.Sprintf(ca.SignClientCertConf(), "CN", "Beijing", "Beijing", "mycompany", "mydept", "myserver")), os.ModePerm)
					execcmd.Command("openssl", "genrsa", "-out", keyPath, "2048")
					execcmd.Command("openssl", "req", "-new", "-sha256", "-out", csrPath, "-key", keyPath, "-config", confFilePath)
					csr, _ := os.ReadFile(csrPath)
					return csr
				}()), WithDays(10)}}, want: "", want1: filepath.Join(homeDir, "/ca/cert/t2/", "client.pem"), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCert(tt.fields.caKey, tt.fields.caPem, tt.fields.opts...)
			got, got1, err := c.signFromCsr()
			if (err != nil) != tt.wantErr {
				t.Errorf("signFromCsr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("signFromCsr() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("signFromCsr() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
