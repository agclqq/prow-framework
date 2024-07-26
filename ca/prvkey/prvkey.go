package prvkey

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func Gen(bit int) ([]byte, error) {
	//创建私钥
	prvKey, err := rsa.GenerateKey(rand.Reader, bit)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(prvKey)}), nil
}
