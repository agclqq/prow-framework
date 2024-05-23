package make

import (
	"fmt"
	"strings"

	file2 "github.com/agclqq/prow-framework/file"
	"github.com/agclqq/prow-framework/prowjob/command"
)

func GetFullPath(defaultDir, sourcePath, fileName string) string {
	if defaultDir != "" && !strings.HasSuffix(defaultDir, "/") {
		defaultDir += "/"
	}
	if sourcePath != "" && !strings.HasSuffix(sourcePath, "/") {
		sourcePath += "/"
	}
	return defaultDir + sourcePath + strings.ToLower(fileName) + ".go"
}

func CheckOverwrite(file string) bool {
	if file2.Exist(file) {
		fmt.Printf(command.FILE_EXIST+" \n", file)
		fmt.Println("whether to overwrite this file?[y/n]")
		goon := ""
		if _, err := fmt.Scanf("%s", &goon); err != nil {
			fmt.Println(err)
			return false
		}
		if goon != "y" {
			return false
		}
	}
	return true
}
