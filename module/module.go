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
		return "", nil
	}
	moduleName = strings.TrimSpace(string(output))
	return moduleName, nil
}
