package execcmd

import (
	"io"
	"os/exec"
)

func Command(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = cmd.Stdout
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	rs, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	return rs, err
}
