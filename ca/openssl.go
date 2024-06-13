package ca

import "github.com/agclqq/prow-framework/execcmd"

func SslVerify() bool {
	log, err := execcmd.Command("openssl", "version")
	if err != nil {
		return false
	}
	if string(log) == "" {
		return false
	}
	return true
}
