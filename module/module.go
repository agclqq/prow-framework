package module

import (
	"os/exec"
	"strings"
)

var moduleName string

// GetName get module name
func GetName() (string, error) {
	if moduleName != "" {
		return moduleName, nil
	}
	cmd := exec.Command("go", "list", "-m")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	moduleName = strings.TrimSpace(string(output))
	return moduleName, nil
}

// GetNameWithoutErr get module name
func GetNameWithoutErr() string {
	name, _ := GetName()
	return name
}

func GetFrameworkName() string {
	return "github.com/agclqq/prow-framework"
}
