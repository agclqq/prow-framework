package codefmt

import (
	"github.com/agclqq/prow-framework/execcmd"
)

func GoModTidy() (string, error) {
	s, err := execcmd.Command("go", "mod", "tidy")
	return string(s), err
}

func GoFmt() (string, error) {
	s, err := execcmd.Command("gofmt", "-w", "-s", ".")
	return string(s), err
}

func GoImports(path ...string) (string, error) {
	if len(path) > 0 {
		s, err := execcmd.Command("goimports", "-local", path[0], "-w", ".")
		return string(s), err
	}
	s, err := execcmd.Command("goimports", "-local", "-w", ".")
	return string(s), err
}
