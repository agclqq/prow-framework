package crl

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"math/big"
)

func InitCrl(caCert *x509.Certificate, caKey crypto.Signer) ([]byte, error) {
	rl := &x509.RevocationList{}

	rl.RevokedCertificateEntries = []x509.RevocationListEntry{}

	rl.Number = big.NewInt(1)

	//rl.ThisUpdate = time.Now()
	//rl.NextUpdate = time.Now().AddDate(0, 0, 7)
	crlBytes, err := x509.CreateRevocationList(rand.Reader, rl, caCert, caKey)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "X509 CRL", Bytes: crlBytes}), nil
}
