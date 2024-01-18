package make

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func Test_createRouter(t *testing.T) {
	t.Skip("此测试函数不会被执行")
	dir, _ := os.Getwd()
	file := dir + "../../../../../application/http/router/router.go"
	fmt.Println(filepath.Abs(file))
	createRouter(file)
}
