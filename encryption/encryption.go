package encryption

import (
	"github.com/agclqq/goencryption"
)

type EasyCrypt struct {
	Type string
	Key  string
	Iv   string
}

func (e *EasyCrypt) Encrypt(plaintext string) (string, error) {
	return goencryption.EasyEncrypt(e.Type, plaintext, e.Key, e.Iv)
}

func (e *EasyCrypt) Decrypt(cipherText string) (string, error) {
	return goencryption.EasyDecrypt(e.Type, cipherText, e.Key, e.Iv)
}
