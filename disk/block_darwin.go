package disk

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"

	"github.com/agclqq/prow-framework/execcmd"
)

func getBlockSize() int {
	out, err := execcmd.Command("diskutil", "info", "/")
	if err != nil {
		return 0
	}
	// 解析输出，查找簇大小
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Block Size:") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				blockSize := parts[2]
				fmt.Println("簇大小:", blockSize)
				return cast.ToInt(blockSize)
			}
		}
	}
	return 0
}
